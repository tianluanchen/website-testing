package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnBadError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, &Error{
		Code:    Status.Bad,
		Message: message,
	})
	ctx.Abort()
}

func ReturnBadErrorWithData(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusBadRequest, &ErrorWithData{
		Universal: Universal{
			Code:    Status.Bad,
			Message: message,
		},
		Data: data,
	})
	ctx.Abort()
}

func ReturnSuccessWithData(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, &UniversalWithData{
		Universal: Universal{
			Code:    Status.OK,
			Message: message,
		},
		Data: data,
	})
	ctx.Abort()
}

func ReturnSuccess(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Universal{
		Code:    Status.OK,
		Message: message,
	})
	ctx.Abort()
}
