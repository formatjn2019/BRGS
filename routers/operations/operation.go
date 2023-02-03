package operations

import (
	"BRGS/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetState(c *gin.Context) {
	t := models.FstreeType{}
	c.JSON(http.StatusOK, map[string]interface{}{
		"state": t.Translate(models.FSTREE_BACKUP),
	})
}

func OperationBackup(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"operation": true,
	})
}
