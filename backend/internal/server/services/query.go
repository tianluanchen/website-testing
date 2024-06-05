package services

import (
	"net/url"
	"website-testing/config"
	"website-testing/internal/server/response"
	"website-testing/internal/tc"

	"github.com/gin-gonic/gin"
)

func GetOptions(ctx *gin.Context) {
	origins := make([]string, len(config.Origins))
	for i := range config.Origins {
		origins[i] = (*url.URL)(config.Origins[i]).String()
	}
	response.ReturnSuccessWithData(ctx, "success", gin.H{
		"dns_servers": config.DNSServerOptions,
		"user_agents": config.UserAgentOptions,
		"origins":     origins,
	})
}

func QueryTestingState(ctx *gin.Context) {
	response.ReturnSuccessWithData(ctx, "success", tc.IsTesting())
}

func GetTestingResult(ctx *gin.Context) {
	response.ReturnSuccessWithData(ctx, "success", tc.GetStore())
}

func GetItemContent(ctx *gin.Context) {
	obj := struct {
		Categroy string `form:"category" binding:"required"`
		Name     string `form:"name" binding:"required"`
	}{}
	if err := ctx.ShouldBindQuery(&obj); err != nil {
		response.ReturnBadError(ctx, err.Error())
		return
	}
	store := tc.GetStore()
	if store == nil {
		response.ReturnBadError(ctx, "还未进行过测试")
		return
	}
	for _, group := range store.Groups {
		if group.Category == obj.Categroy {
			for _, item := range group.Items {
				if item.Name == obj.Name {
					var content []byte
					if item.Result != nil {
						content = item.Result.Content
					}
					response.ReturnSuccessWithData(ctx, "success", string(content))
					return
				}

			}

		}
	}
	response.ReturnBadError(ctx, "查询项不存在")
}
