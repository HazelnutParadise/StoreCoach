package database

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mongoClient *mongo.Client
	MongoDB     *mongo.Database
)

const (
	mongoHost         = "mongo"
	mongoPort         = "27017"
	MongoDatabaseName = "StoreCoach"
)

var NoDBMode = false

func InitMongo(devMode bool) {
	if devMode {
		NoDBMode = true
		return
	}
	var err error
	mongoClient, err = mongo.Connect(options.Client().ApplyURI("mongodb://" + mongoHost + ":" + mongoPort))
	if err != nil {
		log.Fatal(err)
	}
	MongoDB = mongoClient.Database(MongoDatabaseName)
}
