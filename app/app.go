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
	StoreName string   `json:"storeName"`
	Reviews   []string `json:"reviews"`
	timeStamp int64
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

func HandleReviewMining(dataUUID string) (result []SingleReviewMiningResult, err error) {
	reviewData, ok := ReviewMiningDataBuf.LoadAndDelete(dataUUID)
	if !ok {
		err = errors.New("data not found")
		return
	}
	reviews := reviewData.(ReviewData).Reviews
	storeName := reviewData.(ReviewData).StoreName
	result, err = ReviewMining(storeName, reviews)
	return
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
