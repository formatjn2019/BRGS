package routers

import (
	"BRGS/conf"
	"BRGS/routers/operations"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/state", operations.GetState)
	router.GET("/archives", operations.GetArchive)

	router.GET("/zipBackup", operations.OperationZipBackup)
	router.GET("/hardLinkBackup", operations.OperationHardLinkBackup)
	router.GET("/scanning", operations.OperationScanning)
	router.GET("/pause", operations.OperationPause)
	router.GET("/continue", operations.OperationContinue)
	router.GET("/sync", operations.OperationSync)
	router.GET("/reSync", operations.OperationReSync)

	router.POST("/recover", operations.OperationRecover)

	return router
}

func StartServer() error {
	router := InitRouter()
	return router.Run(":" + conf.ServerConf.ServerPort)
}
