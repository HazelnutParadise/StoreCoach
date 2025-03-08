package routes

import (
	"StoreCoach/app"
	"StoreCoach/database"
	"log"
	"net/http"

	_ "embed"

	"github.com/HazelnutParadise/Go-Utils/sliceutil"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine, indexHtml []byte, assets http.FileSystem) {
	apiGp := r.Group("/api")
	apiGp.POST("/review-mining", func(c *gin.Context) {
		// 取得 json 資料，給 uuid，存入 DataBuf
		reviewData := app.ReviewData{}
		err := c.ShouldBindJSON(&reviewData)
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{
				"message": "Bad Request",
			})
			return
		}
		// 清除空的評論
		iters := len(reviewData.Reviews)
		for i := iters - 1; i >= 0; i-- {
			if reviewData.Reviews[i] == "" {
				var err error
				reviewData.Reviews, err = sliceutil.Remove(reviewData.Reviews, i)
				if reviewData.Ratings != nil {
					reviewData.Ratings, err = sliceutil.Remove(reviewData.Ratings, i)
				}
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
		result, _ := database.FindReviewMiningResult(dataUUID)
		if result != nil {
			c.JSON(200, result)
			return
		}
		var err error
		result, err = app.HandleReviewMining(dataUUID)
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
		database.SaveReviewMiningResult(dataUUID, result)
		c.JSON(200, result)
	})

	r.GET("/", func(c *gin.Context) {
		c.Writer.Write(indexHtml)
	})
	r.StaticFS("/assets", assets)
	r.NoRoute(func(c *gin.Context) {
		c.Writer.Write(indexHtml)
	})
}
