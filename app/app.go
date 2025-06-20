package app

import (
	"errors"
	"sync"
	"time"

	"github.com/HazelnutParadise/insyra"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var ReviewMiningDataBuf = sync.Map{}
var IsMiningList = insyra.NewDataList()
var ReviewMiningErrBuf = sync.Map{}

type ReviewData struct {
	StoreName   string   `json:"storeName"`
	ProductName string   `json:"productName"`
	Reviews     []string `json:"reviews"`
	Ratings     []uint8  `json:"ratings"`
	timeStamp   int64
}

func init() {
	go func() {
		for {
			keysToDelete := []string{}
			anHourAgo := time.Now().Add(-time.Hour).Unix()
			ReviewMiningDataBuf.Range(func(key, value any) bool {
				if value.(ReviewData).timeStamp < anHourAgo {
					keysToDelete = append(keysToDelete, key.(string))
				}
				return true
			})
			for _, key := range keysToDelete {
				ReviewMiningDataBuf.Delete(key)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func (reviewData *ReviewData) SetTimeStamp() {
	reviewData.timeStamp = time.Now().Unix()
}

func ReviewMining_SaveToBuf(reviewData ReviewData) (dataUUID string) {
	dataUUID = uuid.New().String()
	reviewData.SetTimeStamp()
	ReviewMiningDataBuf.Store(dataUUID, reviewData)
	return
}

func HandleReviewMining(dataUUID string) (result *ReviewMiningStruct, err error) {
	reviewDataAny, ok := ReviewMiningDataBuf.LoadAndDelete(dataUUID)
	if !ok {
		err = errors.New("data not found")
		ReviewMiningErrBuf.Store(dataUUID, err)
		return nil, err
	}
	IsMiningList.Append(dataUUID)

	reviewData := reviewDataAny.(ReviewData)
	reviews := reviewData.Reviews
	storeName := reviewData.StoreName
	productName := reviewData.ProductName
	ratings := reviewData.Ratings

	result, err = ReviewMining(storeName, productName, reviews, ratings)
	IsMiningList.DropAll(dataUUID)
	return result, err
}

func splitIntoChunks[T any](slice []T, length int) (chunks [][]T) {
	for i := 0; i < len(slice); i += length {
		end := min(i+length, len(slice))
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
