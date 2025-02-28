package app

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

var DataBuf sync.Map
var ResultBuf sync.Map

type ReviewData struct {
	StoreName string   `json:"store_name"`
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
