package Routes

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ViewLogs(GinContext *gin.Context) {

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

	// Get logs from service

	var Logs []Models.Log = Service.Logs

	GinContext.JSON(200, Types.Response{

		Message: fmt.Sprintf("Retrieved %d logs", len(Logs)),

		JSON: Logs,

	});

}