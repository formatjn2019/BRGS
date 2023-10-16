package operations

import (
	"BRGS/management"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetState(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"state":      management.MonServer.State(),
		"reloadFlag": management.MonServer.NeedReload(),
	})
}

func GetArchive(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"archives": management.MonServer.LoadArchives(),
	})
}

func OperationZipBackup(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.ZipArchive(),
	})
}

func OperationHardLinkBackup(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.HardLinkArchive(),
	})
}

func OperationRecover(c *gin.Context) {
	archive := c.PostForm("archiveName")
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.RecoverFormArchive(archive),
	})
}

func OperationPause(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.Pause(),
	})
}

func OperationContinue(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.Continue(),
	})
}

func OperationScanning(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.Scanning(),
	})
}

func OperationReSync(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.Recover(),
	})
}

func OperationSync(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": management.MonServer.Backup(),
	})
}
