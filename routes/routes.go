package routes

import (
	"StoreCoach/app"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	apiGp := r.Group("/api")
	apiGp.POST("/review-mining", func(c *gin.Context) {
		// 取得 json 資料，給 uuid，存入 DataBuf
		reviewData := app.ReviewData{}
		c.BindJSON(&reviewData)
		reviewData.SetTimeStamp()
		dataUUID := app.SaveToBuf(reviewData)
		c.JSON(200, gin.H{
			"data_uuid": dataUUID,
		})
	})

	// TODO
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
