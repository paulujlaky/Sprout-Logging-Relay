package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"elucid503/SproutLoggingRelay/Assets"

	"elucid503/SproutLoggingRelay/Middleware"
	"elucid503/SproutLoggingRelay/Types"
)

var ConfigFilePath string = "./Assets/Config.json"
var DbClient *mongo.Client

func LoadConfig() (*Types.Config, error) {

	// Read the file

	File, Err := os.Open(ConfigFilePath)

	if Err != nil {

		return nil, Err

	}

	defer File.Close()

	// Decode the file

	Decoder := json.NewDecoder(File)

	Config := &Types.Config{} // Get the pointer to the Config struct

	if Err := Decoder.Decode(Config); Err != nil {

		return nil, Err

	}

	return Config, nil

}

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

var RouteHandlers = Assets.GetRoutes()

func LoadRoutes(Router *gin.Engine) {

	for Route, Handler := range RouteHandlers {

		MiddlewareFuncs := []gin.HandlerFunc{}

		if Handler.Authorized {

			MiddlewareFuncs = append(MiddlewareFuncs, Middleware.Authorize())

		}

		if Handler.RateLimitConfig != nil {

			MiddlewareFuncs = append(MiddlewareFuncs, Middleware.RateLimit(Handler.RateLimitConfig))

		}

		MiddlewareFuncs = append(MiddlewareFuncs, Handler.Handler)

		Router.Handle(Handler.Method, Route, MiddlewareFuncs...)

	}

}

func main() {

	Config, Error := LoadConfig()

	if Error != nil {

		panic(Error) // Can't continue

	}

	fmt.Printf("Sprout Logging Relay Backend; Version %s\n", Config.Versions.Server)

	gin.SetMode(Config.GinMode)

	GinRouter := gin.Default()

	LoadRoutes(GinRouter)

	fmt.Printf("Loaded %d routes...\n", len(RouteHandlers))

	// Connect to the database

	DbError := ConnectToDb(Config)

	if DbError != nil {

		panic(DbError) // Can't continue

	}

	fmt.Println("Connected to the database...")

	fmt.Printf("Listening on port %d...\n", Config.Server.Port)

	// Serve Static files

	GinRouter.Static("/Out", "../Frontend/Out")
	GinRouter.Static("/Assets", "../Frontend/Assets")
	GinRouter.Static("/Styles", "../Frontend/Styles") // for maps back to SCSS

	GinRouter.Run(fmt.Sprintf(":%d", Config.Server.Port))

}
