package Routes

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func GetServiceList(GinContext *gin.Context) {

	// Get all services

	Services, ErrorGettingServices := Models.ListServices()

	if ErrorGettingServices != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Error getting services",

		})

		return

	}

	GinContext.JSON(200, Types.Response{

		Message: "Retrieved services",

		JSON: Services,

	})
	
}
