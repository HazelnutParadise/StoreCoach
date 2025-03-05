package app

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var ReviewMiningDataBuf sync.Map

type ReviewData struct {
	StoreName   string   `json:"storeName"`
	ProductName string   `json:"productName"`
	Reviews     []string `json:"reviews"`
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
	reviewData, ok := ReviewMiningDataBuf.LoadAndDelete(dataUUID)
	if !ok {
		err = errors.New("data not found")
		return nil, err
	}
	reviews := reviewData.(ReviewData).Reviews
	storeName := reviewData.(ReviewData).StoreName
	productName := reviewData.(ReviewData).ProductName
	result, err = ReviewMining(storeName, productName, reviews)
	return result, err
}

func splitIntoChunks[T any](slice []T, length int) (chunks [][]T) {
	for i := 0; i < len(slice); i += length {
		end := i + length
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
