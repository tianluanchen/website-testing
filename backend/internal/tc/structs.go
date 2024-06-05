package tc

import (
	"sync"
	"time"
	"website-testing/config"
	"website-testing/wt"
)

type (
	addInfo struct {
		// tcp or udp
		Type string `json:"type"`
		IP   string `json:"ip"`
		Port int    `json:"port"`
	}
	respInfo struct {
		StatusCode    int    `json:"status_code"`
		Status        string `json:"status"`
		ContentLength int64  `json:"content_length"`
		ContentType   string `json:"content_type"`
	}
	recordDetail struct {
		URL string `json:"url"`
		// maybe nil
		RemoteAddr *addInfo `json:"remote_addr"`
		// milliseconds
		Duration int64     `json:"duration"`
		Resp     *respInfo `json:"resp"`
	}
	testResult struct {
		// milliseconds
		TotalDuration int64           `json:"total_duration"`
		Records       []*recordDetail `json:"records"`
		// content size
		Size  int    `json:"size"`
		Title string `json:"title"`
		// ignore
		Content []byte `json:"-"`
		Err     string `json:"err,omitempty"`
		// if the last response is a redirect, assign redirect url
		LastRespRedirect string `json:"last_resp_redirect,omitempty"`
	}
	testItem struct {
		URLWithName
		// pending or done
		Status string      `json:"status"`
		Result *testResult `json:"result"`
		Err    string      `json:"err,omitempty"`
	}
	testGroup struct {
		Category string      `json:"category"`
		Items    []*testItem `json:"items"`
	}
	testingStore struct {
		Conf   *config.Conf `json:"conf"`
		Groups []*testGroup `json:"groups"`
		// unix milliseconds
		Start int64 `json:"start"`
		// unix milliseconds
		End int64  `json:"end"`
		Err string `json:"err,omitempty"`
	}
)

type testingCenter struct {
	client  *wt.Client
	store   *testingStore
	testing bool
	rwmutex sync.RWMutex
}

type CallbackOption struct {
	OnStart         func()
	OnRaceForAPI    func(api *config.API, duration time.Duration)
	OnFetchWebsites func(count int)
	OnTest          func(count, finished int, category, name, link string)
	OnFinish        func(err error, duration time.Duration)
}

func (opt *CallbackOption) fix() {
	if opt == nil {
		return
	}
	if opt.OnStart == nil {
		opt.OnStart = func() {}
	}
	if opt.OnRaceForAPI == nil {
		opt.OnRaceForAPI = func(api *config.API, duration time.Duration) {}
	}
	if opt.OnFetchWebsites == nil {
		opt.OnFetchWebsites = func(count int) {}
	}
	if opt.OnTest == nil {
		opt.OnTest = func(count, finished int, category, name, link string) {}
	}
	if opt.OnFinish == nil {
		opt.OnFinish = func(err error, duration time.Duration) {}
	}
}
