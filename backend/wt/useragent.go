package wt

import (
	"math/rand"
	"strings"
	"time"
)

const (
	UAChromeWin64   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	UAChromeMacOS   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	UAChromeLinux   = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	UAChromeiOS     = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/122.0.6261.89 Mobile/15E148 Safari/604.1"
	UAChormeAndroid = "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.119 Mobile Safari/537.36"

	UASafariMacOS = "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3.1 Safari/605.1.15"
	UASafariOS    = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3.1 Mobile/15E148 Safari/604.1"

	UAFirefoxWin64   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"
	UAFirefoxMacOS   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 14.4; rv:115.0) Gecko/20100101 Firefox/115.0"
	UAFirefoxLinux   = "Mozilla/5.0 (X11; Linux i686; rv:123.0) Gecko/20100101 Firefox/123.0"
	UAFirefoxiOS     = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/123.0 Mobile/15E148 Safari/605.1.15"
	UAFirefoxAndroid = "Mozilla/5.0 (Android 14; Mobile; rv:123.0) Gecko/123.0 Firefox/123.0"
)

var UserAgentMap = map[string]string{
	"ChromeWin64":   UAChromeWin64,
	"ChromeMacOS":   UAChromeMacOS,
	"ChromeLinux":   UAChromeLinux,
	"ChromeiOS":     UAChromeiOS,
	"ChormeAndroid": UAChormeAndroid,

	"SafariMacOS": UASafariMacOS,
	"SafariOS":    UASafariOS,

	"FirefoxWin64":   UAFirefoxWin64,
	"FirefoxMacOS":   UAFirefoxMacOS,
	"FirefoxLinux":   UAFirefoxLinux,
	"FirefoxiOS":     UAFirefoxiOS,
	"FirefoxAndroid": UAFirefoxAndroid,
}

// if substr is empty, return a random user agent
//
// if specify one substr at least, return a  random user agent that contains all substr(Case-insensitive), if neither contains, return UAChromeWin64
func GetUserAgent(substr ...string) string {
	for i := 0; i < len(substr); i++ {
		substr[i] = strings.ToLower(substr[i])
	}
	list := []string{}

	for k, v := range UserAgentMap {
		k := strings.ToLower(k)
		valid := true
		if len(substr) > 0 {
			for _, s := range substr {
				if !strings.Contains(k, s) {
					valid = false
					break
				}
			}
		}
		if valid {
			list = append(list, v)
		}
	}
	if len(list) == 0 {
		return UAChromeWin64
	}
	return list[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(list))]
}
