package Models

type LogLevel int

const (

	// Log Levels

	LogLevelInfo    LogLevel = 0
	LogLevelWarning LogLevel = 1
	LogLevelError   LogLevel = 2
)

type LogFeatureType string

const (

	// Log Feature Types

	LogFeatureTypeCode LogFeatureType = "code"
	LogFeatureTypeLink LogFeatureType = "link"
)

func IsValidFeatureType(FeatureType string) bool {

	// See if string is within the defined feature types

	switch FeatureType {

	case string(LogFeatureTypeCode):

		return true

	case string(LogFeatureTypeLink):

		return true

	default:

		return false

	}

}

type LogFeature struct {


	// May represent code blocks, links, or anything else the parser may need to identify
	// The renderer will handle the rendering of these respective features


	Type  LogFeatureType `json:"Type"`
	Value string         `json:"Value"`

}

type LogContent struct {

	Title       string `json:"Title"`
	Description string `json:"Description"`

	Features []LogFeature `json:"Features"`

}

type Log struct {

	Level LogLevel `json:"Level"`

	Content LogContent `json:"Content"`

	Created int64 `json:"Created"`
	
}
