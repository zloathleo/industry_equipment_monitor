package realtime

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRealtimeController(dasGroup *gin.RouterGroup) {

	dasGroup.GET("/realtime", func(c *gin.Context) {
		points := c.Query("names")

		bytes, err := generatePointRealtimeValue(points)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			return
		}
		c.Status(http.StatusOK)
		c.Writer.Write(bytes)
	})

	dasGroup.POST("/push", func(c *gin.Context) {
		dataJson := c.PostForm("content")
		_, err := pushRealtimeValue(dataJson)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
			})
		}

	})

}
