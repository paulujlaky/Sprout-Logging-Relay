package Models

import (
	"context"
	"strings"
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
	
	Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").InsertOne(context.TODO(), Service)

	return Service

}

func ListServices() ([]Service, error) {

	var Services []Service = []Service{}

	Options := options.Find().SetProjection(gin.H{"logs": 0}) // exclude logs from the response since it would be too large and not make sense in this context
	
	DBCursor, Err := Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").Find(context.TODO(), gin.H{}, Options) // get all services

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

	Err := Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").FindOne(context.TODO(), gin.H{"uid": UID}).Decode(&ServiceData)

	if Err != nil {

		return nil, Err
		
	}

	return &ServiceData, nil
	
}

// Internal Funcs

func (AssociatedService *Service) parseRawLogContent(RawContent string) []LogFeature {

	// features will be denoted by {{feature_type:feature_value}}

	var ParsedFeatures []LogFeature = []LogFeature{}

	// these will be updated throughout

	SplitString := strings.Split(RawContent, "{{")

	for _, Split := range SplitString {

		if strings.Contains(Split, "}}") {

			// we have a feature, now we just need to parse it

			FeatureSplit := strings.Split(Split, "}}")

			FeatureTypeValue := strings.Split(FeatureSplit[0], ":")

			if len(FeatureTypeValue) == 2 && IsValidFeatureType(FeatureTypeValue[0]) {

				Feature := LogFeature{

					Type:  LogFeatureType(FeatureTypeValue[0]),
					Value: FeatureTypeValue[1],
				}

				ParsedFeatures = append(ParsedFeatures, Feature)

			}

		}

	}

	return ParsedFeatures

}

// Public Funcs

func (AssociatedService *Service) AddLog(Level LogLevel, Title string, RawContent string) error {
	
	if (AssociatedService.MaxLogs < len(AssociatedService.Logs)) {

		// we need to remove the oldest log

		AssociatedService.Logs = AssociatedService.Logs[1:] // reassign slice to everything except the first element

	}

	AssociatedService.Logs = append(AssociatedService.Logs, Log{

		Level: Level,

		Content: LogContent{

			Title:       Title,
			Description: RawContent,
			Features:    AssociatedService.parseRawLogContent(RawContent),
		},

		Created: time.Now().Unix(),
	})

	// Update in DB 

	_, WriteError := Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").UpdateOne(context.TODO(), gin.H{"uid": AssociatedService.UID}, gin.H{"$set": gin.H{"logs": AssociatedService.Logs}})

	return WriteError

}