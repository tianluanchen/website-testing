package tc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"website-testing/config"
	"website-testing/pkg"
	"website-testing/wt"

	"golang.org/x/sync/errgroup"
)

var logger = pkg.NewLogger()

// Unique instance
var tc = &testingCenter{}

func (tc *testingCenter) buildClient(conf *config.Conf) error {
	client, err := wt.New(conf.ToWtOption())
	if err != nil {
		return err
	}
	tc.client = client
	return nil
}

func (tc *testingCenter) raceForAPI() (*config.API, *wt.Result, error) {
	var winner *config.API
	var winnerResult *wt.Result
	var once sync.Once
	var wg sync.WaitGroup
	errs := make([]error, len(config.Origins))
	ch := make(chan struct{})
	for i, v := range config.Origins {
		i, u := i, (*url.URL)(v).String()
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := tc.client.Visit(u)
			if err != nil {
				errs[i] = err
				return
			}
			if result.Err != nil {
				errs[i] = result.Err
				return
			}
			if len(result.Records) > 1 {
				errs[i] = errors.New("exist redirects and last response is from " + result.Records[len(result.Records)-1].Request.URL.String())
				return
			}
			once.Do(func() {
				ch <- struct{}{}
				winner = v
				winnerResult = result
			})
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	<-ch
	if winner == nil {
		reason := "all APIs encountered execution errors:\n"
		for i, v := range errs {
			reason += fmt.Sprintf("Error using API %s: %s\n", (*url.URL)(config.Origins[i]).String(), v.Error())
		}
		return nil, nil, errors.New(reason)
	}
	return winner, winnerResult, nil
}

type URLWithName struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (tc *testingCenter) FetchWebsites(api *config.API) (map[string][]URLWithName, error) {
	result := make(map[string][]URLWithName)
	g := new(errgroup.Group)
	for _, v := range []string{config.CategoryAnimation, config.CategoryVideo} {
		u := v
		g.Go(func() (fatalErr error) {
			defer func() {
				if fatalErr != nil {
					logger.Errorf("捕获分类%s的网站列表失败: %s", u, fatalErr.Error())
				}
			}()
			logger.Debugln("开始捕获分类" + u + "的网站列表...")
			req, err := http.NewRequest(http.MethodGet, api.GetURLWithCategory(u), nil)
			if err != nil {
				return err
			}
			records, err := tc.client.Do(req, 0)
			if err != nil {
				return err
			}

			list, err := parseResponse(records[len(records)-1].Response)
			if err != nil {
				return err
			}
			logger.Infof("成功从API(%s)捕获分类%s的网站列表", (*url.URL)(api).String(), u)
			result[u] = list
			return nil
		})
	}
	return result, g.Wait()
}

func (tc *testingCenter) IsTesting() bool {
	tc.rwmutex.RLock()
	defer tc.rwmutex.RUnlock()
	return tc.testing
}
func (tc *testingCenter) Test(conf *config.Conf, opt *CallbackOption) bool {
	if tc.IsTesting() {
		return false
	}
	tc.rwmutex.Lock()
	tc.testing = true
	tc.rwmutex.Unlock()
	go func() {
		defer func() {
			tc.rwmutex.Lock()
			tc.testing = false
			tc.rwmutex.Unlock()
		}()
		tc.test(conf, opt)
	}()
	return true
}

func getErrValue(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
func getTimeValue(t time.Time) int64 {
	return t.UnixMilli()
}
func getDurationValue(d time.Duration) int64 {
	return d.Milliseconds()
}
func getAddrInfo(addr net.Addr) *addInfo {
	if addr == nil {
		return nil
	}
	switch v := addr.(type) {
	case *net.TCPAddr:
		return &addInfo{
			"tcp",
			v.IP.String(),
			v.Port,
		}
	case *net.UDPAddr:
		return &addInfo{
			"udp",
			v.IP.String(),
			v.Port,
		}
	default:
		return nil
	}
}

func getRespInfo(resp *http.Response) *respInfo {
	if resp == nil {
		return nil
	}
	return &respInfo{
		resp.StatusCode,
		resp.Status,
		resp.ContentLength,
		resp.Header.Get("Content-Type"),
	}
}

var titleRegexp = regexp.MustCompile(`<title[^<>]*>([^<>]*)</title>`)

func parseTitleFromHTML(content []byte) string {
	result := titleRegexp.FindSubmatch(content)
	if result == nil {
		return ""
	}
	return string(result[1])
}

func parseVisitReturn(result *wt.Result, err error) (*testResult, string) {
	var r *testResult
	if result != nil {
		r = &testResult{
			TotalDuration: getDurationValue(result.TotalDuration),
			Title:         parseTitleFromHTML(result.Content),
			Size:          len(result.Content),
			Content:       result.Content,
			Err:           getErrValue(result.Err),
		}
		for i, v := range result.Records {
			r.Records = append(r.Records, &recordDetail{
				URL:        v.Request.URL.String(),
				RemoteAddr: getAddrInfo(v.RemoteAddr),
				Duration:   getDurationValue(v.Duration),
				Resp:       getRespInfo(v.Response),
			})
			if i == len(result.Records)-1 && v.Response != nil {
				if u, err := v.Response.Location(); err == nil {
					r.LastRespRedirect = u.String()
				}
			}

		}

	}
	return r, getErrValue(err)

}
func (tc *testingCenter) test(conf *config.Conf, opt *CallbackOption) (fatalErr error) {
	if opt == nil {
		opt = &CallbackOption{}
	}
	opt.fix()
	start := time.Now()
	tc.store = &testingStore{
		Conf:  conf,
		Start: getTimeValue(start),
	}
	defer func() {
		end := time.Now()
		tc.store.End = getTimeValue(end)
		tc.store.Err = getErrValue(fatalErr)
		opt.OnFinish(fatalErr, end.Sub(start))
	}()
	opt.OnStart()
	logger.Debugln("配置项:", conf)
	if err := tc.buildClient(conf); err != nil {
		logger.Errorln("构建客户端失败:", err)
		return err
	}
	logger.Debugln("正在选择最快的可用API...")
	api, result, err := tc.raceForAPI()
	if err != nil {
		logger.Errorln("获取可用API失败:", err)
		return err
	}

	opt.OnRaceForAPI(api, result.TotalDuration)

	logger.Debugln("选用API:", (*url.URL)(api).String(), "访问时长:", result.TotalDuration)
	websites, err := tc.FetchWebsites(api)
	if err != nil {
		logger.Errorln("获取网站失败:", err)
		return err
	}

	count := 0
	groups := make([]*testGroup, 0)
	for c, list := range websites {
		group := &testGroup{
			Category: c,
			Items:    make([]*testItem, 0),
		}
		for _, v := range list {
			count++
			group.Items = append(group.Items, &testItem{
				URLWithName: v,
				Status:      "pending",
			})
		}
		groups = append(groups, group)
	}
	opt.OnFetchWebsites(count)
	tc.store.Conf = conf
	tc.store.Groups = groups
	concurrency := count
	if count > 64 {
		concurrency = 64
	}
	ch := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var finished atomic.Int64
	logger.Debugln("开始测试", count, "个网站...")
	for i := range groups {
		for j := range groups[i].Items {
			item := groups[i].Items[j]
			wg.Add(1)
			ch <- struct{}{}
			go func() {
				defer func() {
					wg.Done()
					<-ch
				}()
				item.Result, item.Err = parseVisitReturn(tc.client.Visit(item.URL))
				item.Status = "done"
				n := finished.Add(1)
				opt.OnTest(count, int(n), groups[i].Category, groups[i].Items[j].Name, groups[i].Items[j].URL)
			}()
		}
	}
	wg.Wait()
	logger.Debug("测试完成")
	return nil
}

func parseResponse(resp *http.Response) ([]URLWithName, error) {
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("invalid status: " + resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "json") {
		return nil, errors.New("invalid content type: " + contentType)
	}
	obj := struct {
		Websites []URLWithName `json:"websites"`
	}{}
	encoder := json.NewDecoder(resp.Body)
	if err := encoder.Decode(&obj); err != nil {
		return nil, err
	}
	return obj.Websites, nil
}

func IsTesting() bool {
	return tc.IsTesting()
}

func Test(conf *config.Conf, opt *CallbackOption) bool {
	return tc.Test(conf, opt)
}

// maybe nil and the struct has custom marshaller
func GetStore() *testingStore {
	return tc.store
}
