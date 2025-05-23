package Models

import (
	"context"
	"fmt"
	"time"

	"elucid503/SproutLoggingRelay/Functions"
	"elucid503/SproutLoggingRelay/Misc"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service struct {

	UID string `json:"UID"`

	Name        string `json:"Name"`
	Description string `json:"Description"`
	Icon        string `json:"Icon"`

	// Internals 

	MaxLogs int `json:"MaxLogs"` // max logs to keep before overwriting

}

// Creator 

func NewService(Name string, Description string, Icon string) Service {

	// Generate UID

	UID := Functions.RandomString(16)

	Service := Service{

		UID:         UID,
		Name:        Name,
		Description: Description,
		Icon:        Icon,

		MaxLogs:     250, // default max logs to keep before overwriting. eventually might be configurable by user

	}
	
	Misc.Database.Collection("Services").InsertOne(context.TODO(), Service)

	return Service

}

func ListServices() ([]Service, error) {

	var Services []Service

	Cursor, Err := Misc.Database.Collection("Services").Find(context.TODO(), gin.H{})

	if Err != nil {

		return nil, Err
		
	}

	for Cursor.Next(context.TODO()) {

		var Service Service

		Err := Cursor.Decode(&Service)

		if Err != nil {

			return nil, Err
			
		}

		Services = append(Services, Service)

	}

	return Services, nil

}

func ListLogsFromDB(ServiceUID string, Offset int, Limit int) ([]*Log, error) {

	Filter := gin.H{}

	if (ServiceUID != "All") {

		// We want specific logs from a service

		Filter = gin.H{"parentserviceuid": ServiceUID}	

	}

	Options := options.Find()

	Options.SetSort(gin.H{"created": -1}) // sort by created date (oldest first)
	Options.SetSkip(int64(Offset))       // skip the first Offset logs
	Options.SetLimit(int64(Limit))       // limit the number of logs to Limit
	Options.SetProjection(gin.H{"_id": 0}) // exclude the _id field from the results

	Cursor, Err := Misc.Database.Collection("Logs").Find(context.TODO(), Filter, Options)

	if Err != nil {

		return nil, Err
		
	}

	var Logs []*Log

	for Cursor.Next(context.TODO()) {

		var LogData Log

		Err := Cursor.Decode(&LogData)

		if Err != nil {

			return nil, Err
			
		}

		Logs = append(Logs, &LogData)

	}

	// Reverse the logs to show the most recent first

	for i, j := 0, len(Logs)-1; i < j; i, j = i+1, j-1 {

		Logs[i], Logs[j] = Logs[j], Logs[i]

	}

	return Logs, nil

}

// Getter (From DB)

func GetServiceByUID(UID string) (*Service, error) {

	var ServiceData Service

	Err := Misc.Database.Collection("Services").FindOne(context.TODO(), gin.H{"uid": UID}).Decode(&ServiceData)

	if Err != nil {

		return nil, Err
		
	}

	return &ServiceData, nil
	
}

// Public Funcs

func (AssociatedService *Service) AddLog(Level LogLevel, Title string, RawContent string) (*Log, error) {

	// RULES:
	// Logs must be stored in the Logs collection 
	// Logs for each service cannot exceed Service.MaxLogs 

	NewLog := Log{

		UID: Functions.RandomString(16), 
		ParentServiceUID: AssociatedService.UID, 

		Level: Level,

		Content: LogContent{

			Title:       Title,
			Description: RawContent,
		},

		Created: time.Now().Unix(),

	}

	NewLog.Content.Features = NewLog.parseRawLogContent() // parse the raw log content into features
	NewLog.HTML = NewLog.ToHTML(AssociatedService) // convert the log to HTML

	NumberOfLogs, CntError := Misc.Database.Collection("Logs").CountDocuments(context.TODO(), gin.H{"ParentServiceUID": AssociatedService.UID})

	if CntError != nil {

		fmt.Println("Error counting logs:", CntError)

		return nil, CntError

	}

	if NumberOfLogs >= int64(AssociatedService.MaxLogs) {

		// delete the oldest log

		_, DelError := Misc.Database.Collection("Logs").DeleteOne(context.TODO(), gin.H{"ParentServiceUID": AssociatedService.UID})

		if DelError != nil {

			fmt.Println("Error deleting log:", DelError)

			return nil, DelError

		}

	}
	
	_, InsertErr := Misc.Database.Collection("Logs").InsertOne(context.TODO(), NewLog) // insert the log into the database

	if InsertErr != nil {

		fmt.Println("Error inserting log:", InsertErr)

		return nil, InsertErr
		
	}

	return &NewLog, nil

}