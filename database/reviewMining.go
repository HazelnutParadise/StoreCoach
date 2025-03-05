package database

import (
	"StoreCoach/app"
	"encoding/json"
	"errors"
)

type ReviewMiningResult struct {
	DataUUID   string `gorm:"primaryKey"`
	DataString string
}

func FindReviewMiningResult(dataUUID string) (*app.ReviewMiningStruct, error) {
	var result = ReviewMiningResult{}
	rows := DB.Where("data_uuid = ?", dataUUID).Find(&result).RowsAffected
	if rows == 0 {
		return nil, errors.New("data not found")
	}
	var resultStruct app.ReviewMiningStruct
	err := json.Unmarshal([]byte(result.DataString), &result)
	if err != nil {
		return nil, err
	}
	return &resultStruct, nil
}
func SaveReviewMiningResult(dataUUID string, result *app.ReviewMiningStruct) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}
	var obj = ReviewMiningResult{
		DataUUID:   dataUUID,
		DataString: string(jsonData),
	}
	DB.Save(obj)
	return nil
}
