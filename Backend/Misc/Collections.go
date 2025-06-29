package Misc

import (
	"context"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/gorilla/websocket"
)

// Globals

var DbClient *mongo.Client
var Database *mongo.Database

var ConnectedWebSockets []*websocket.Conn

var Config *Types.Config

// Database

func ConnectToDb(Config *Types.Config) error {

	ClientOptions := options.Client().ApplyURI(Config.Database.URI)

	Client, Err := mongo.Connect(nil, ClientOptions)

	if Err != nil {

		return Err

	}

	// Assign the client to the global variable

	DbClient = Client
	Database = Client.Database(Config.Database.Name)

	return nil

}

// WebSockets

func AddWebSocket(Conn *websocket.Conn) {

	ConnectedWebSockets = append(ConnectedWebSockets, Conn)

}

func RemoveWebSocket(Conn *websocket.Conn) {

	for i, WebSocket := range ConnectedWebSockets {

		if WebSocket == Conn {

			ConnectedWebSockets = slices.Delete(ConnectedWebSockets, i, i+1)

			break

		}

	}

}

func SendToAllWebSockets(Message Types.Response) {

	for _, WebSocket := range ConnectedWebSockets {

		WebSocket.WriteJSON(Message)

	}

}

func WatchAndDeleteOldLogsFromDatabase() {

	// This will create a loop which runs every 5 minutes, and deletes logs older than 5 days from the database

	for {

		CurrentTime := time.Now()
		FiveDaysAgo := CurrentTime.AddDate(0, 0, -5).Unix()

		_, Err := Database.Collection("Logs").DeleteMany(context.Background(), bson.M{"created": bson.M{"$lt": FiveDaysAgo}})

		if Err != nil {

			fmt.Println("[CLEANUP] Error deleting old logs from database:", Err)

		} else {

			fmt.Println("[CLEANUP] Deleted old logs from database")

		}

		time.Sleep(5 * time.Minute)

	}

}