package httpr

import (
	"github.com/gin-gonic/gin"
)

// 访问域名限制
func newCrosMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	}
}

// access_token
//func newAuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		access_token := c.Request.Header.Get("access_token")
//		if access_token == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"error": "no access token",
//			})
//			c.Abort()
//			return
//		} else {
//			c.Set("access_token", access_token)
//			c.Next()
//			return
//		}
//	}
//}
