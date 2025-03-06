package database

import (
	"StoreCoach/app"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func FindReviewMiningResult(dataUUID string) (*app.ReviewMiningStruct, error) {
	var result app.ReviewMiningStruct
	found := MongoDB.Collection("reviewMiningResult").FindOne(context.Background(), bson.M{"dataUUID": dataUUID})
	if err := found.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
func SaveReviewMiningResult(dataUUID string, result *app.ReviewMiningStruct) error {
	result.DataUUID = dataUUID
	_, err := MongoDB.Collection("reviewMiningResult").InsertOne(context.Background(), result)
	return err
}
