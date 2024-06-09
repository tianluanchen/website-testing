package wt

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"runtime"
	"time"
)

func getProxyFunc(proxyURL string) (func(*http.Request) (*url.URL, error), error) {
	if proxyURL == "" {
		return http.ProxyFromEnvironment, nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "socks5" && u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported scheme " + u.Scheme)
	}
	return http.ProxyURL(u), nil
}

func getDialContextWithSpecifiedDNS(dnsServer string) (func(ctx context.Context, network string, address string) (net.Conn, error), error) {
	if dnsServer == "" {
		return nil, nil
	}
	// if _, err := net.ResolveTCPAddr("tcp", dnsServer); err != nil {
	// 	return nil, errors.New("invalid dns server")
	// }
	var dnsServerAddr string
	addr, err := net.ResolveIPAddr("ip", dnsServer)
	if err != nil {
		return nil, err
	}
	if ip := addr.IP.To4(); ip != nil {
		dnsServerAddr = ip.String() + ":53"
	} else if ip := addr.IP.To16(); ip != nil {
		dnsServerAddr = "[" + ip.String() + "]:53"
	} else {
		return nil, errors.New("invalid dns server")
	}
	d := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := &net.Dialer{}
				return d.DialContext(ctx, network, dnsServerAddr)
			},
		},
	}
	return d.DialContext, nil
}

func traceReq(req *http.Request, trace *httptrace.ClientTrace) *http.Request {

	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}

// EnsureTransporterFinalized will ensure that when the HTTP client is GCed
// the runtime will close the idle connections (so that they won't leak)
// this function was adopted from Hashicorp's go-cleanhttp package
func ensureTransporterFinalized(httpTransport *http.Transport) {
	runtime.SetFinalizer(&httpTransport, func(transportInt **http.Transport) {
		(*transportInt).CloseIdleConnections()
	})
}
