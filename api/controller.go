package api

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/raafvargas/wrapit/contract"
)

// Controller ...
type Controller interface {
	RegisterRoutes(*gin.RouterGroup)
}

//BindingError ...
func BindingError(context *gin.Context, err error) {
	context.Error(err)
	context.JSON(http.StatusBadRequest, contract.NewError(http.StatusBadRequest, err.Error()))
}

// ResolveError ...
func ResolveError(ctx *gin.Context, err error) {
	ctx.Error(err)

	if reflect.TypeOf(err) != reflect.TypeOf(&contract.Error{}) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	status := http.StatusBadRequest
	message := err.(*contract.Error)

	if message.Code != 0 {
		status = message.Code
	}

	ctx.AbortWithStatusJSON(status, message)
}
