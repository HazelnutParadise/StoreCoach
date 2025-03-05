package routes

import (
	"StoreCoach/app"

	"github.com/HazelnutParadise/Go-Utils/sliceutil"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	apiGp := r.Group("/api")
	apiGp.POST("/review-mining", func(c *gin.Context) {
		// 取得 json 資料，給 uuid，存入 DataBuf
		reviewData := app.ReviewData{}
		c.BindJSON(&reviewData)
		// 清除空的評論
		iters := len(reviewData.Reviews)
		for i := iters - 1; i >= 0; i-- {
			if reviewData.Reviews[i] == "" {
				var err error
				reviewData.Reviews, err = sliceutil.Remove(reviewData.Reviews, i)
				if err != nil {
					c.JSON(400, gin.H{
						"message": "Bad Request",
					})
					return
				}
			}
		}
		dataUUID := app.ReviewMining_SaveToBuf(reviewData)
		c.JSON(200, gin.H{
			"dataUUID": dataUUID,
		})
	})
	apiGp.GET("/review-mining/:data_uuid", func(c *gin.Context) {
		dataUUID := c.Param("data_uuid")
		result, err := app.HandleReviewMining(dataUUID)
		if err != nil {
			if err.Error() == "data not found" {
				c.JSON(404, gin.H{
					"message": "Data not found",
				})
				return
			}
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
			})
		}
		c.JSON(200, result)
	})

	// TODO
	r.GET("/vite.svg", func(c *gin.Context) {
		c.File("frontend/dist/vite.svg")
	})
	r.Static("/assets", "frontend/dist/assets")
	r.NoRoute(func(c *gin.Context) {
		c.File("frontend/dist/index.html")
	})
}
