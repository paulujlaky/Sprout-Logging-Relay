package Routes

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func Add(GinContext *gin.Context) {

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	Subscription := Body["Subscription"]
	Key := Body["Key"]
	PublicKey := Body["PublicKey"]

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

		Details: Subscription.(gin.H),
		
		Key: Key.(string),
		PublicKey: PublicKey.(string),

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