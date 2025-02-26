package app

import (
	"sync"

	"github.com/google/uuid"
)

var DataBuf sync.Map
var ResultBuf sync.Map

type ReviewData struct {
	StoreName string   `json:"store_name"`
	Reviews   []string `json:"reviews"`
}

func SaveToBuf(reviewData ReviewData) (dataUUID string) {
	dataUUID = uuid.New().String()
	DataBuf.Store(dataUUID, reviewData)
	return
}
