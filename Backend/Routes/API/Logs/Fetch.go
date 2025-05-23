package Routes

import (
	"elucid503/SproutLoggingRelay/Misc"
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func FetchLog(GinContext *gin.Context) {

	LogUID := GinContext.Param("LogUID")

	if LogUID == "" {	
		
		GinContext.JSON(400, Types.Response{

			Message: "Missing Log UID",

		})

		return

	}

	// We will try to get the log from the DB by UID

	LogFromDatabase := Models.Log{}

	ErrorFetching := Misc.Database.Collection("Logs").FindOne(GinContext, gin.H{"UID": LogUID}).Decode(&LogFromDatabase)

	if ErrorFetching != nil || LogFromDatabase.UID == "" {

		GinContext.JSON(404, Types.Response{

			Message: "Log not found",

		})

		return

	}
	
	GinContext.JSON(200, Types.Response{

		Message: "Log fetched successfully",
		JSON: LogFromDatabase,

	})

}