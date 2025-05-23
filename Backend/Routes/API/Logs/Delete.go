package Routes

import (
	"elucid503/SproutLoggingRelay/Misc"
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func DeleteLog(GinContext *gin.Context) {

	LogUID := GinContext.Param("LogUID")

	if LogUID == "" {	
		
		GinContext.JSON(400, Types.Response{

			Message: "Missing Log UID",

			JSON: gin.H{

				"Success": false,

			},

		})

		return

	}

	// We will try to delete the log from the DB by UID

	DeleteResult := Misc.Database.Collection("Logs").FindOneAndDelete(GinContext, gin.H{"uid": LogUID})

	if DeleteResult.Err() != nil {

		GinContext.JSON(404, Types.Response{

			Message: "Log could not be deleted",

			JSON: gin.H{

				"Success": false,

			},

		})

		return

	}
	
	GinContext.JSON(200, Types.Response{

		Message: "Log deleted successfully",

		JSON: gin.H{

			"Success": true,

		},
		
	})

}