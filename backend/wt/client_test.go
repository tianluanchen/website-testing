package wt

import (
	"testing"
)

func printResult(t *testing.T, result *Result) {
	for i, r := range result.Records {
		var status string
		if r.Response != nil {
			status = r.Response.Status
		}
		t.Log("url:", r.Request.URL.String(), "remote addr:", r.RemoteAddr, "status:", status, "duration:", r.Duration)
		if i == len(result.Records)-1 {
			if u, err := r.Response.Location(); err == nil {
				t.Log("the final response want to redirect to:", u.String())
			}
		}
	}
	t.Log("total duration:", result.TotalDuration, "err:", result.Err)
	t.Log("content:\n", string(result.Content))
}

func TestClient(t *testing.T) {
	client, err := New(&Option{
		UserAgent: UAFirefoxWin64,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Visit("https://httpbin.org/anything")
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, result)
}

func TestClientWithRedirects(t *testing.T) {
	client, err := New(&Option{
		UserAgent: UAFirefoxWin64,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 2 redirects
	result, err := client.Visit("http://httpbin.org/redirect-to?url=http%3A%2F%2Fhttpbin.org%2Fredirect-to%3Furl%3Dhttps%253A%252F%252Fhttpbin.org%252Fanything")
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, result)
}

func TestClientWithLimitRedirects(t *testing.T) {
	client, err := New(&Option{
		UserAgent: UAFirefoxWin64,
	})
	if err != nil {
		t.Fatal(err)
	}
	// maximum 1 redirect, but the link will redirect twice
	result, err := client.Visit(
		"http://httpbin.org/redirect-to?url=http%3A%2F%2Fhttpbin.org%2Fredirect-to%3Furl%3Dhttps%253A%252F%252Fhttpbin.org%252Fanything",
		1,
	)
	if err != nil {
		t.Fatal(err)
	}
	printResult(t, result)
}
