package services

import (
	"net/url"
	"time"
	"website-testing/config"
	"website-testing/internal/server/response"
	"website-testing/internal/tc"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

var wsHandler = melody.New()

func init() {
	wsHandler.HandleConnect(func(s *melody.Session) {
		s.Write(serializeWithPanic(gin.H{
			"event": "State",
			"data":  tc.IsTesting(),
		}))
	})
	wsHandler.HandleMessage(func(s *melody.Session, msg []byte) {
		if len(msg) == 4 && string(msg) == "ping" {
			s.Write([]byte("pong"))
		}
	})
}

var CallbackOption = &tc.CallbackOption{
	OnStart: func() {
		wsHandler.Broadcast(serializeWithPanic(gin.H{
			"event": "Start",
		}))
	},
	OnPickFastestAPI: func(api *config.API, duration time.Duration) {
		wsHandler.Broadcast(serializeWithPanic(gin.H{
			"event": "PickFastestAPI",
			"data": gin.H{
				"api":      (*url.URL)(api).String(),
				"duration": duration.Milliseconds(),
			},
		}))
	},
	OnFetchWebsites: func(count int) {
		wsHandler.Broadcast(serializeWithPanic(gin.H{
			"event": "FetchWebsites",
			"data":  count,
		}))
	},
	OnTest: func(count, finished int, category, name, link string) {
		wsHandler.Broadcast(serializeWithPanic(gin.H{
			"event": "Test",
			"data": gin.H{
				"category": category,
				"name":     name,
				"count":    count,
				"finished": finished,
				"url":      link,
			},
		}))
	},
	OnFinish: func(err error, duration time.Duration) {
		obj := gin.H{
			"duration": duration.Milliseconds(),
			"err":      nil,
		}
		if err != nil {
			obj["err"] = err.Error()
		}
		wsHandler.Broadcast(serializeWithPanic(gin.H{
			"event": "Finish",
			"data":  obj,
		}))
	},
}

func StartTesting(ctx *gin.Context) {
	conf := &config.Conf{}
	if err := ctx.ShouldBindJSON(conf); err != nil {
		response.ReturnBadError(ctx, err.Error())
		return
	}
	response.ReturnSuccessWithData(ctx, "success", gin.H{
		"last_test_aborted": tc.Test(conf, CallbackOption),
	})
}

func WatchTestingStatus(ctx *gin.Context) {
	wsHandler.HandleRequest(ctx.Writer, ctx.Request)
}
