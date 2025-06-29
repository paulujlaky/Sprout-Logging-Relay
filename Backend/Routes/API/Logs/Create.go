package Routes

import (
	"elucid503/SproutLoggingRelay/Misc"
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateLog(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	// Check body

	if (Body["Level"] == nil) || (Body["Title"] == nil) || (Body["Content"] == nil) {

		GinContext.JSON(400, Types.Response{

			Message: "Missing required fields",

		})

	}

	// Parse level as int

	ParsedLevel, ErrorParsingLevel := strconv.Atoi(fmt.Sprintf("%v", Body["Level"]))

	if ErrorParsingLevel != nil {

		GinContext.JSON(400, Types.Response{

			Message: "Level is not a number",

		})

		return

	}

	// Check if Level is valid

	if !Models.IsValidLogLevel(ParsedLevel) {

		GinContext.JSON(400, Types.Response{

			Message: "Level is invalid",

		})

		return

	}	
	
	// Parse body

	Title := Body["Title"].(string)
	Content := Body["Content"].(string)

	// Get UID from request URL params 

	UID := GinContext.Param("ServiceUID")

	if UID == "" {

		GinContext.JSON(400, Types.Response{

			Message: "No UID provided",

		})
			
		return

	}

	// Get service 

	Service, ErrorGettingService := Models.GetServiceByUID(UID)

	if ErrorGettingService != nil {

		GinContext.JSON(404, Types.Response{

			Message: "Service not found",

		})

		return

	}

	// Create log

	NewLog, AddLogError := Service.AddLog(Models.LogLevel(ParsedLevel), Title, Content)
	
	if AddLogError != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Error adding log",

		})

	} else {

		GinContext.JSON(200, Types.Response{

			Message: "Log added successfully",

		})
		
		// In this case, we also should send to all websockets waiting for logs that this log was created 
		
		Misc.SendToAllWebSockets(Types.Response{
			
			Message: "New log created",
			JSON: gin.H{
				"ServiceUID": Service.UID,
				"Log": NewLog,
			},
		})

		// Now we must check if it was 
		// a. an error
		// b. if there was more than n errors in the last m seconds (as determined by the linked services's notification settings)
		// this takes some database reads... so we will only check it if it is an error log 

		if ParsedLevel == int(Models.LogLevelError) {

			// good to go 

			if Service.Notifications.Threshold <= 0 || Service.Notifications.Interval <= 0 { return }

			// Get the number of errors in the last n seconds, including this one 

			NumberOfErrors, CountError := Misc.Database.Collection("Logs").CountDocuments(GinContext, gin.H{

				"parentserviceuid": Service.UID,
				"level":            Models.LogLevelError,

				"created": gin.H{

					"$gte": NewLog.Created - int64(Service.Notifications.Interval),

				},

			})

			if CountError != nil { return } // not ideal.

			if NumberOfErrors >= int64(Service.Notifications.Threshold) {

				// the service is currently not doing good 

				AllNotificationSubscriptions, ErrorGettingSubscriptions := Models.GetAllSubscriptions();

				if ErrorGettingSubscriptions == nil { return }
				
				for _, Subscription := range AllNotificationSubscriptions {

					if (time.Now().Unix() - Subscription.LastNotifSent) < 300 {
						
						return // we have already sent a notification in the last 5 minutes, no need to send another one

					}

					// Send notification to subscription



					

				}



			}


			


				


		}

	}

}