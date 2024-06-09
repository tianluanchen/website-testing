package server

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"website-testing/internal/web"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	var skipPaths []string
	fs.WalkDir(web.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if !strings.HasSuffix(d.Name(), ".html") {
				skipPaths = append(skipPaths, "/"+path)
			}
		}
		return nil
	})
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	fileServer := http.FileServer(http.FS(web.FS))
	etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(web.MD5Hash))
	app.Use(func(ctx *gin.Context) {
		if _, ok := web.PathMap[ctx.Request.URL.Path]; ok {
			ctx.Header("ETag", etag)
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}
	})
	app.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skipPaths,
	}))
	initRouter(app)
	return app
}
