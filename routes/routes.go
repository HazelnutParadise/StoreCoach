package routes

import "github.com/gin-gonic/gin"

func SetRoutes(r *gin.Engine) {
	// Set routes here
	r.POST("/", func(c *gin.Context) {
		// TODO: 取得 json 資料，給 uuid，存入 DataBuf
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
