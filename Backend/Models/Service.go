package Models

import (
	"context"
	"sync"
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

	Logs []Log `json:"Logs"`

	LogsMutex sync.Mutex // mutex for thread-safe log updates. NOT to be used for DB updates, and NOT serialized to JSON

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

		Logs:        []Log{},

		MaxLogs:     500, // default max logs to keep before overwriting. eventually might be configurable by user

	}
	
	Misc.Database.Collection("Services").InsertOne(context.TODO(), Service)

	return Service

}

func ListServices(WithLogs bool) ([]Service, error) {

	var Services []Service = []Service{}

	Options := options.Find()

	if !WithLogs {

		Options.SetProjection(gin.H{"logs": 0}) // exclude logs from the response using this projection
	
	} 
		
	DBCursor, Err := Misc.Database.Collection("Services").Find(context.TODO(), gin.H{}, Options) // get all services

	if Err != nil {

		return nil, Err

	}

	DBCursor.All(context.TODO(), &Services) // unmarshal all entries into Services slice

	DBCursor.Close(context.TODO()) // after reading, free memory

	return Services, nil
	
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

func (AssociatedService *Service) AddLog(Level LogLevel, Title string, RawContent string) (error, *Log) {
	
	AssociatedService.LogsMutex.Lock() // lock the mutex for thread-safe in-memory updates
	defer AssociatedService.LogsMutex.Unlock()

	if (AssociatedService.MaxLogs < len(AssociatedService.Logs)) {

		// we need to remove the oldest log

		AssociatedService.Logs = AssociatedService.Logs[1:] // reassign slice to everything except the first element

	}

	NewLog := Log{

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

	// Update in-memory logs

	AssociatedService.Logs = append(AssociatedService.Logs, NewLog)

	if len(AssociatedService.Logs) > AssociatedService.MaxLogs {

		AssociatedService.Logs = AssociatedService.Logs[1:] // trim oldest log

	}

	// Update in DB atomically
	
	_, WriteError := Misc.Database.Collection("Services").UpdateOne(
		context.TODO(),
		gin.H{"uid": AssociatedService.UID},
		gin.H{
			"$push": gin.H{"logs": gin.H{
				"$each": []Log{NewLog},
				"$slice": -AssociatedService.MaxLogs, // ensure logs array doesn't exceed MaxLogs
			}},
		},
	)

	return WriteError, &NewLog // return the new log and any errors that occurred

}