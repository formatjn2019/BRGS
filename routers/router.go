package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {

	})

	return r
}
