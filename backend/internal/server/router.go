package server

import (
	"website-testing/internal/server/services"

	"github.com/gin-gonic/gin"
)

func initRouter(app *gin.Engine) {
	g := app.Group("/api")
	g.GET("/options", services.GetOptions)
	g.GET("/state", services.QueryTestingState)
	g.POST("/start", services.StartTesting)
	g.GET("/result", services.GetTestingResult)
	g.GET("/content", services.GetItemContent)
	g.GET("/watch", services.WatchTestingStatus)
}
