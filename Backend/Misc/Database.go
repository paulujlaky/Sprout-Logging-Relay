package Misc

import (
	"elucid503/SproutLoggingRelay/Types"
	"slices"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/gorilla/websocket"
)

// Globals

var DbClient *mongo.Client

var ConnectedWebSockets []*websocket.Conn

// Database

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