package wt

import (
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"time"
)

type Client struct {
	transport *http.Transport
	ua        string
	timeout   time.Duration
}

type Option struct {
	UserAgent string
	// if empty or invalid, use http.ProxyFromEnvironment
	ProxyURL string
	// hostname
	DNSServer string
	// if zero, no timeout
	Timeout time.Duration
}

func (client *Client) ChangeUA(ua string) {
	client.ua = ua
}

func New(option *Option) (*Client, error) {
	if option == nil {
		option = &Option{}
	}
	client := &Client{
		timeout: option.Timeout,
	}
	var err error
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// unlimit
	transport.MaxConnsPerHost = 0
	transport.MaxIdleConns = 0
	transport.MaxIdleConnsPerHost = 10000

	transport.DialContext, err = getDialContextWithSpecifiedDNS(option.DNSServer)
	if err != nil {
		return nil, err
	}
	transport.Proxy, err = getProxyFunc(option.ProxyURL)
	if err != nil {
		return nil, err
	}
	client.transport = transport
	ensureTransporterFinalized(client.transport)
	client.ChangeUA(option.UserAgent)
	return client, nil
}

// if error is nil, you should close last response.Body manually
//
// default maximumRedirects is 10
func (client *Client) Do(req *http.Request, maximumRedirects ...int) ([]Record, error) {
	if len(maximumRedirects) == 0 {
		maximumRedirects = []int{10}
	}
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", client.ua)
	}
	jar, _ := cookiejar.New(nil)

	records := make([]Record, 0)
	remoteAddrList := make([]net.Addr, 0)
	respList := make([]*http.Response, 0)
	reqList := []*http.Request{req}
	durationList := make([]time.Duration, 0)
	var start time.Time

	req = traceReq(req, &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddrList = append(remoteAddrList, connInfo.Conn.RemoteAddr())
		},
	})

	c := &http.Client{
		Transport: client.transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > maximumRedirects[0] {
				return http.ErrUseLastResponse
			}
			durationList = append(durationList, time.Since(start))
			respList = append(respList, req.Response)
			reqList = append(reqList, req)
			// update
			start = time.Now()
			return nil
		},
		Jar:     jar,
		Timeout: client.timeout,
	}

	start = time.Now()

	resp, err := c.Do(req)
	if err == nil {
		respList = append(respList, resp)
		durationList = append(durationList, time.Now().Sub(start))
	}
	for i, req := range reqList {
		record := Record{
			Request: req,
		}
		if len(remoteAddrList) > i {
			record.RemoteAddr = remoteAddrList[i]
		}
		if len(respList) > i {
			record.Response = respList[i]
		}
		if len(durationList) > i {
			record.Duration = durationList[i]
		}
		records = append(records, record)
	}
	return records, err
}

// visit url like a browser with GET method
//
// if err is not nil, it means first request build error
func (client *Client) Visit(url string, maximumRedirects ...int) (*Result, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language": []string{"zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2"},
		"Sec-Fetch-Dest":  []string{"document"},
		"Sec-Fetch-Mode":  []string{"navigate"},
		"Sec-Fetch-Site":  []string{"none"},
		"Sec-Fetch-User":  []string{"?1"},
	}
	result := &Result{}
	start := time.Now()
	result.Records, result.Err = client.Do(req, maximumRedirects...)
	// it means we can't get final Response
	if result.Err != nil {
		result.TotalDuration = time.Since(start)
		return result, err
	}
	resp := result.Records[len(result.Records)-1].Response
	defer resp.Body.Close()
	// Maximum read 512KB
	result.Content, result.Err = io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	result.TotalDuration = time.Since(start)
	return result, nil
}

type Record struct {
	Request *http.Request
	// maybe nil
	Response *http.Response
	// maybe nil
	RemoteAddr net.Addr
	// maybe zero
	Duration time.Duration
}

type Result struct {
	// maybe nil
	Records []Record
	// maybe nil
	Content       []byte
	TotalDuration time.Duration
	// maybe is last response body read error
	Err error
}
