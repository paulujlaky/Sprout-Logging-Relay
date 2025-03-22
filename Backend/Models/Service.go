package Models

import (
	"context"
	"strings"
	"time"

	"elucid503/SproutLoggingRelay/Misc"
)

type Service struct {

	UID string `json:"UID"`

	Name        string `json:"Name"`
	Description string `json:"Description"`
	Icon        string `json:"Icon"`

	Logs []Log `json:"Logs"`

}

// Creator 

func NewService(UID string, Name string, Description string, Icon string) Service {

	Service := Service{

		UID:         UID,
		Name:        Name,
		Description: Description,
		Icon:        Icon,

		Logs:        []Log{},

	}
	
	Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").InsertOne(context.TODO(), Service)

	return Service

}

// Getter (From DB)

func GetServiceByUID(UID string) (*Service, error) {

	var ServiceData Service

	Err := Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").FindOne(context.TODO(), Service{UID: UID}).Decode(&ServiceData)

	if Err != nil {

		return nil, Err
		
	}

	return &ServiceData, nil
	
}

// Internal Funcs

func (AssociatedService *Service) parseRawContent(RawContent string) []LogFeature {

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

func (AssociatedService *Service) AddLog(Level LogLevel, Title string, RawContent string) {

	AssociatedService.Logs = append(AssociatedService.Logs, Log{

		Level: Level,

		Content: LogContent{

			Title:       Title,
			Description: RawContent,
			Features:    AssociatedService.parseRawContent(RawContent),
		},

		Created: time.Now().Unix(),
	})

	// Update in DB 

	Misc.DbClient.Database("SproutLoggingRelay").Collection("Services").UpdateOne(context.TODO(), Service{UID: AssociatedService.UID}, AssociatedService)

}
