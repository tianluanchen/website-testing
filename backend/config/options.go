package config

import (
	"sort"
	"website-testing/wt"
)

type option[T string] struct {
	Label string `json:"label"`
	Value T      `json:"value"`
}

var DNSServerOptions = []option[string]{
	{
		Label: "跟随系统",
		Value: "",
	},
	{
		Label: "腾讯云",
		Value: "119.29.29.29",
	},
	{
		Label: "阿里云",
		Value: "223.5.5.5",
	},
	{
		Label: "Cloudflare",
		Value: "1.1.1.1",
	},
}

var UserAgentOptions = []option[string]{}

func init() {
	for k, v := range wt.UserAgentMap {
		UserAgentOptions = append(UserAgentOptions, option[string]{
			Label: k,
			Value: v,
		})
	}
	sort.Slice(UserAgentOptions, func(i, j int) bool {
		return UserAgentOptions[i].Label < UserAgentOptions[j].Label
	})
}
