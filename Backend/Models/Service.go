package Models

import "time"

type Service struct {

	Name        string `json:"Name"`
	Description string `json:"Description"`
	Icon        string `json:"Icon"`

	Logs []Log `json:"Logs"`
	
}

func (AssociatedService *Service) AddLog(Level LogLevel, Title string, RawContent string) {

	AssociatedService.Logs = append(AssociatedService.Logs, Log{

		Level: Level,

		Content: LogContent{

			Title:       Title,
			Description: RawContent,
			Features:    []LogFeature{},

		},

		Created: time.Now().Unix(),
		
	})

}