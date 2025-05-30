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
		log.Printf("Received review data: %+v", reviewData)
		dataUUID := app.ReviewMining_SaveToBuf(reviewData)
		c.JSON(200, gin.H{
			"dataUUID": dataUUID,
		})
	})
	apiGp.GET("/review-mining/:data_uuid", func(c *gin.Context) {
		dataUUID := c.Param("data_uuid")
		log.Printf("Received request for dataUUID: %s", dataUUID)

		// 先檢查資料庫是否有結果
		result, err := database.FindReviewMiningResult(dataUUID)
		if result != nil {
			c.JSON(200, result)
			return
		} else if err == nil {
			// 資料存在但還沒處理完成，回傳 nil
			c.JSON(200, nil)
			return
		}

		// 非同步開始處理，立即回傳 nil
		go func() {
			log.Printf("Starting async processing for dataUUID: %s", dataUUID)
			result, err := app.HandleReviewMining(dataUUID)
			if err != nil {
				log.Printf("Error processing dataUUID %s: %v", dataUUID, err)
			} else {
				database.SaveReviewMiningResult(dataUUID, result)
				log.Printf("Completed processing for dataUUID: %s", dataUUID)
			}
		}()

		// 立即回傳 nil，表示處理已開始
		c.JSON(200, nil)
	})

	r.GET("/", func(c *gin.Context) {
		c.Writer.Write(indexHtml)
	})
	r.StaticFS("/assets", assets)
	r.NoRoute(func(c *gin.Context) {
		c.Writer.Write(indexHtml)
		c.Data(200, "text/html", indexHtml)
	})
}
