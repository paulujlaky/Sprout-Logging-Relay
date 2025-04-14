package Routes

import (
	"elucid503/SproutLoggingRelay/Functions"
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ViewServiceSpecificLogs(GinContext *gin.Context) {

	// Get service 

	ServiceUID := GinContext.Param("ServiceUID")

	if ServiceUID == "" {

		GinContext.JSON(400, Types.Response{

			Message: "Service UID is required",

			JSON: nil,

		})

		return

	}

	Service, ServiceFetchErr := Models.GetServiceByUID(ServiceUID) // TODO: Migrate this to caching, if we're happy with the performance gains

	if ServiceFetchErr != nil {

		GinContext.JSON(404, Types.Response{

			Message: fmt.Sprintf("Service with UID %s not found", ServiceUID),

			JSON: nil,

		})

		return

	}

	AllLogsHTMLArr := make([]string, len(Service.Logs))

	for i, Log := range Service.Logs {

		AllLogsHTMLArr[i] = Log.HTML

	}

	GinContext.JSON(200, Types.Response{

		Message: fmt.Sprintf("Successfully retrieved %d %s", len(Service.Logs), Functions.PluralizeString(len(Service.Logs), "log")),
		JSON: AllLogsHTMLArr,

	});

}