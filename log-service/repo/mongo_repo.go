package repo

import (
	"context"
	"log"
	"log_service/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo interface {
	InsertData(req *entity.LogEntry) error
}

type MongoRepoImpl struct {
	Collection *mongo.Collection
}

func NewMongoRepo(collection *mongo.Collection) MongoRepo {
	return &MongoRepoImpl{
		Collection: collection,
	}
}

func (repo *MongoRepoImpl) InsertData(req *entity.LogEntry) error {
	_, err := repo.Collection.InsertOne(context.TODO(), req)
	if err != nil {
		log.Println("error inserting data : ", err)
		return err
	}

	return nil
}
