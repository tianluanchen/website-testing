package config

import "net/url"

const (
	CategoryAnimation = "animation"
	CategoryVideo     = "video"
)

type API url.URL

// get api url clone
func (api *API) Unwrap() *url.URL {
	clone := *api
	return (*url.URL)(&(clone))
}

// use const CategoryAnimation or CategoryVideo
func (api *API) GetURLWithCategory(category string) string {
	u := api.Unwrap()
	switch category {
	case CategoryAnimation:
		u.Path = "/ayouth/json/animation.json"
	case CategoryVideo:
		u.Path = "/ayouth/json/video.json"
	default:
		panic("invalid category")
	}
	return u.String()
}

var Origins = make([]*API, 0)

func init() {
	for _, v := range []string{
		"https://ayouth.top",
		"https://ayouth.eu.org",
	} {
		u, err := url.Parse(v)
		if err != nil {
			panic(err)
		}
		if u.Hostname() == "" || u.Scheme != "https" && u.Scheme != "http" {
			panic("invalid url")
		}
		u.RawPath = ""
		u.Path = ""
		Origins = append(Origins, (*API)(u))
	}
	if len(Origins) == 0 {
		panic("no origins")
	}
}
