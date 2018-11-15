package history

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitHistoryController(dasGroup *gin.RouterGroup) {

	dasGroup.GET("/history", func(c *gin.Context) {
		points := c.Query("points")
		to := c.Query("to")   //结束时间
		dur := c.Query("dur") //时长
		interval := c.Query("interval")
		chartType := c.Query("type")

		bytes, err := generatePointHistoryChart(chartType, points, to, dur, interval)
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

}
