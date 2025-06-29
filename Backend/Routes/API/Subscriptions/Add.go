package Routes

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
)

func Add(GinContext *gin.Context) {

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	Subscription := Body["Subscription"]

	if Subscription == nil {

		GinContext.JSON(400, Types.Response{

			Message: "Missing required field: Subscription",

			JSON: gin.H{ 

				"Success": false,

			},

		})

		return

	}

	SubscriptionInstance := Models.Subscription{

		LastNotifSent: 0, // nothing sent yet

		Subscription: Subscription.(webpush.Subscription),
		
	}

	Err := SubscriptionInstance.Save()

	if Err != nil {

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