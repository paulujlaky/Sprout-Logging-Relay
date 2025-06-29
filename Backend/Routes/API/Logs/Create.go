package Routes

import (
	"elucid503/SproutLoggingRelay/Functions"
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

		return;

	}

	GinContext.JSON(200, Types.Response{

		Message: "Log added successfully",

	})
	
	// Extra stuff which runs concurrently

	go func() {

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

		if ParsedLevel != int(Models.LogLevelError) { return } // not an error
		if Service.Notifications.Threshold <= 0 || Service.Notifications.Interval <= 0 { return } // not configured to notify

		// Get the number of errors in the last n seconds, including this one 

		NumberOfErrors, CountError := Misc.Database.Collection("Logs").CountDocuments(GinContext, gin.H{

			"parentserviceuid": Service.UID,
			"level":            Models.LogLevelError,

			"created": gin.H{

				"$gte": NewLog.Created - int64(Service.Notifications.Interval),

			},

		})

		if CountError != nil { return } // not ideal.
		if (NumberOfErrors) < int64(Service.Notifications.Threshold) { return } // not enough

		// the service is currently not doing good 

		AllNotificationSubscriptions, ErrorGettingSubscriptions := Models.GetAllSubscriptions();
		if ErrorGettingSubscriptions != nil { return }
		
		for _, Subscription := range AllNotificationSubscriptions {

			if (time.Now().Unix() - Subscription.LastNotifSent) < int64(Misc.Config.WebPush.NotificationTimeout) { return }

			// Send notification to subscription

			ErrorSending := Subscription.SendNotification(Service.Name, fmt.Sprintf("Potential issue detected, with %d %s in the last %d seconds. Open for more details.", NumberOfErrors, Functions.PluralizeString(int(NumberOfErrors), "error"), Service.Notifications.Interval), false, fmt.Sprintf("https://relay.sprout.software?Service=%s", Service.UID))

			if ErrorSending != nil { fmt.Println("Error sending notification:", ErrorSending) }

		}
			
	}()

}