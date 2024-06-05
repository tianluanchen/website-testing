package config

import (
	"fmt"
	"time"
	"website-testing/wt"
)

// all client config
type Conf struct {
	DNSServer      string `json:"dns_server" binding:"omitempty,ip"`
	ProxyURL       string `json:"proxy_url" binding:"omitempty,url"`
	TimeoutSeconds int    `json:"timeout_seconds" binding:"required,min=1"`
	UserAgent      string `json:"user_agent"`
}

func (conf *Conf) String() string {
	if conf == nil {
		return ""
	}
	return fmt.Sprintf("DNSServer: %s, ProxyURL: %s, TimeoutSeconds: %d, UserAgent: %s",
		conf.DNSServer, conf.ProxyURL, conf.TimeoutSeconds, conf.UserAgent)
}

func (conf *Conf) ToWtOption() *wt.Option {
	if conf == nil {
		return nil
	}
	return &wt.Option{
		UserAgent: conf.UserAgent,
		ProxyURL:  conf.ProxyURL,
		DNSServer: conf.DNSServer,
		Timeout:   time.Duration(conf.TimeoutSeconds) * time.Second,
	}
}
