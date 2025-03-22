package Misc

import (
	"elucid503/SproutLoggingRelay/Types"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DbClient *mongo.Client

func ConnectToDb(Config *Types.Config) error {

	ClientOptions := options.Client().ApplyURI(Config.Database.URI)

	Client, Err := mongo.Connect(nil, ClientOptions)

	if Err != nil {

		return Err

	}

	// Assign the client to the global variable

	DbClient = Client

	return nil

}