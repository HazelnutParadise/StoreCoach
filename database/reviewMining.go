package database

import (
	"StoreCoach/app"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func FindReviewMiningResult(dataUUID string) (*app.ReviewMiningStruct, error) {
	if NoDBMode {
		return nil, nil
	}
	var result app.ReviewMiningStruct
	found := MongoDB.Collection("reviewMiningResult").FindOne(context.Background(), bson.M{"dataUUID": dataUUID})
	if found.Err() != nil {
		if app.IsMiningList.FindFirst(dataUUID) != nil {
			// If the data is still being mined, return nil
			return nil, nil
		}
		return nil, found.Err()
	}
	if err := found.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
func SaveReviewMiningResult(dataUUID string, result *app.ReviewMiningStruct) error {
	if NoDBMode {
		return nil
	}
	if result == nil {
		return errors.New("result is nil")
	}
	result.DataUUID = dataUUID
	_, err := MongoDB.Collection("reviewMiningResult").InsertOne(context.Background(), result)
	return err
}
