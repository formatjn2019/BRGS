package routers

import (
	"BRGS/routers/operations"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		type Button struct {
			Text string
			Url  string
		}

		c.HTML(http.StatusOK, "temple.html", gin.H{
			"title": "备份文件",
			// "single": "Main website",
			"data": []Button{
				{Text: "状态", Url: "state"},
				{Text: "日志", Url: "logs"},
				{Text: "备份", Url: "operation"},
				// {Text: "按钮3", Url: "test"},
				// {Text: "按钮4", Url: "test"},
				// {Text: "按钮5", Url: "test2"},
				// {Text: "按钮6", Url: "test2"},
				// {Text: "按钮7", Url: "test2"},
			},
		})
	})

	router.GET("/state", operations.GetState)
	router.POST("/operation", operations.OperationBackup)
	return router
}
