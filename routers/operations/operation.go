package operations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetState(c *gin.Context) {
	c.JSON(http.StatusOK,
		"stats")
}
