package app

import (
	"sync"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var DataBuf sync.Map
var ResultBuf sync.Map

type ReviewData struct {
	StoreName string   `json:"storeName"`
	Reviews   []string `json:"reviews"`
	timeStamp int64
}

func (reviewData *ReviewData) SetTimeStamp() {
	reviewData.timeStamp = time.Now().Unix()
}

func SaveToBuf(reviewData ReviewData) (dataUUID string) {
	dataUUID = uuid.New().String()
	DataBuf.Store(dataUUID, reviewData)
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
