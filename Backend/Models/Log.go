package Models

import (
	"fmt"
	"strings"
	"time"
)

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

func IsValidLogLevel(Level int) bool {

	// See if string is within the defined log levels

	switch Level {

	case int(LogLevelInfo):

		return true

	case int(LogLevelWarning):

		return true

	case int(LogLevelError):

		return true

	default:

		return false

	}

}

func LogLevelToString(Level LogLevel) string {

	// Convert log level to string

	switch Level {

	case LogLevelInfo:

		return "Info"

	case LogLevelWarning:

		return "Warning"

	case LogLevelError:

		return "Error"

	default:

		return "Unknown"

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

	ParentServiceUID string `json:"ParentServiceUID"`

	Level LogLevel `json:"Level"`

	Content LogContent `json:"Content"`

	Created int64 `json:"Created"`

	HTML string `json:"HTML"`

}

// Private Methods 

func (Log *Log) parseRawLogContent() []LogFeature {

	// features will be denoted by {{feature_type:feature_value}}

	var ParsedFeatures []LogFeature = []LogFeature{}

	// these will be updated throughout

	SplitString := strings.Split(Log.Content.Description, "{{")
	CleanedContent := ""

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
				
				CleanedContent += FeatureSplit[1] // Append the remaining content after the feature

			} else {

				CleanedContent += "{{" + Split // If not a valid feature, keep the original text

			}

		} else {

			// If no feature, keep the original text

			CleanedContent += Split

		}

	}

	// Update the raw content to exclude parsed features

	Log.Content.Description = CleanedContent

	return ParsedFeatures
	
}

// Public Methods 

func (Log *Log) ToHTML(RelatedService *Service) string {

	// Convert log to HTML string

	StringLevel := LogLevelToString(Log.Level)

	LogContent := Log.Content.Description;

	for _, Feature := range Log.Content.Features {

		switch Feature.Type {

			case LogFeatureTypeCode: {

				// Render code block

				LogContent += fmt.Sprintf(`

					<div class="InlineFeature Code">

						<pre><code>%s</code></pre>

					</div>`, 
					
				Feature.Value)

			}

			case LogFeatureTypeLink: {

				// Render link

				LogContent += fmt.Sprintf(`

					<div class="InlineFeature Link">

						<a href="%s" target="_blank">%s</a>

					</div>`,
			
				Feature.Value, Feature.Value)

			}

		}
		
	}

	return fmt.Sprintf(`

	<div class="InlineItem Log %s" ServiceUID="%s">

		<div class="LogMeta">

			<div class="LogMetaGroup">

				<div class="LogMetaService">

					%s

				</div>

				<span class="LogMetaSeparator">•</span>

				<div class="LogMetaLevel">

					%s

				</div>

			</div>

			<div class="LogMetaGroup">

				<div class="LogMetaTimestamp">

					<div class="LogMetaTimestampDate">

						%s

					</div>

					<span class="LogMetaSeparator">•</span>

					<div class="LogMetaTimestampTime">

						%s

					</div>
				
				</div>

			</div>

		</div>

		<div class="LogTitle">

			%s

		</div>

		<div class="LogContent">

			%s

		</div>

	</div>`, StringLevel, RelatedService.UID, RelatedService.Name, StringLevel, 
	
	time.Unix(Log.Created, 0).Format("Jan 2, 2006"), 
	time.Unix(Log.Created, 0).Format("3:04 PM"),

	Log.Content.Title, 
	Log.Content.Description);

}