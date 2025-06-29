package Routes

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
)

func Add(GinContext *gin.Context) {

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	Subscription := Body["Subscription"] // string JSON of webpush.Subscription

	if Subscription == nil {

		GinContext.JSON(400, Types.Response{

			Message: "Missing required field: Subscription",

			JSON: gin.H{ 

				"Success": false,

			},

		})

		return

	}

	// we must convert the Subscription to a webpush.Subscription type

	ConvertedSubscription := webpush.Subscription{}

	ParsingErr := json.Unmarshal([]byte(Subscription.(string)), &ConvertedSubscription)

	if ParsingErr != nil {

		GinContext.JSON(400, Types.Response{

			Message: "Invalid Subscription format",
			JSON: gin.H{

				"Success": false,

			},

		})

	}

	SubscriptionInstance := Models.Subscription{

		LastNotifSent: 0, // nothing sent yet

		Subscription: ConvertedSubscription, // convert to webpush.Subscription type
		
	}

	SaveErr := SubscriptionInstance.Save()

	if SaveErr != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Failed to save subscription",

			JSON: gin.H{

				"Success": false,

			},

		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Subscription added successfully",

		JSON: gin.H{

			"Success": true,

		},
		
	})

}