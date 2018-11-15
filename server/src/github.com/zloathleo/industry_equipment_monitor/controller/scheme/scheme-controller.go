package scheme

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitSchemeController(dasGroup *gin.RouterGroup) {

	dasGroup.GET("/scheme", func(c *gin.Context) {
		name := c.Query("name")
		var bytes []byte
		var err error
		if name == "" {
			//全查
			bytes, err = generatePointSchemeMap()
		} else {
			//指定查找
			bytes, err = generatePointScheme(name)
		}
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
