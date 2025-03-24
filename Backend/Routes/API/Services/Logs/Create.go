package Routes

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"
	"strconv"

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

	AddLogError := Service.AddLog(Models.LogLevel(ParsedLevel), Title, Content)
	
	if AddLogError != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Error adding log",

		})

		return

	} else {

		GinContext.JSON(200, Types.Response{

			Message: "Log added successfully",

		})

	}

}