package Routes

import (
	"elucid503/SproutLoggingRelay/Functions"
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewServiceSpecificLogs(GinContext *gin.Context) {

	StrOffset := GinContext.Query("Offset")
	StrLimit := GinContext.Query("Limit")

	IntOffset, _ := strconv.Atoi(StrOffset)
	IntLimit, _ := strconv.Atoi(StrLimit)

	ServiceUID := GinContext.Param("ServiceUID")

	if ServiceUID == "" {

		GinContext.JSON(400, Types.Response{

			Message: "Service UID is required",

			JSON: nil,

		})

		return

	}

	ServiceLogs, ServiceLogsFetchErr := Models.ListLogsFromDB(ServiceUID, IntOffset, IntLimit)

	if ServiceLogsFetchErr != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Error fetching logs from database",
			JSON: nil,

		})

	}

	GinContext.JSON(200, Types.Response{

		Message: fmt.Sprintf("Successfully retrieved %d %s", len(ServiceLogs), Functions.PluralizeString(len(ServiceLogs), "log")),
		JSON: ServiceLogs,

	});

}