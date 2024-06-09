package web

import (
	"crypto/md5"
	"embed"
	"io"
	"io/fs"
	"strings"
)

//go:embed all:dist/*
var distFS embed.FS

// static resource
var FS fs.FS

// Map of the url path
var PathMap = map[string]string{}

var MD5Hash []byte

func init() {
	var err error
	FS, err = fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	// url root path
	if ExistFile("index.html") {
		PathMap["/"] = "index.html"
	}
	hasher := md5.New()
	fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		urlPath := "/" + path
		if d.IsDir() {
			urlPath := urlPath + "/"
			indexPage := path + "/index.html"
			if ExistFile(indexPage) {
				PathMap[urlPath] = indexPage
			}
		} else {
			PathMap[urlPath] = path
			f, _ := FS.Open(path)
			io.Copy(hasher, f)
		}
		return nil
	})
	MD5Hash = hasher.Sum(nil)
}

// Determine if a file exists in a path
func ExistFile(path string) bool {
	path = strings.TrimLeft(path, "/")
	f, err := FS.Open(path)
	if err != nil {
		return false
	}
	info, _ := f.Stat()
	return !info.IsDir()
}
