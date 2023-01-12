package connection

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnection() (*mongo.Client, error) {
	const uri = "mongodb://mongo:27017"

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting to mongo: ", err)
		return nil, err
	}

	return c, nil
}
