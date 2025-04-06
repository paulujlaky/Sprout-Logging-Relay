package Routes

import (
	"elucid503/SproutLoggingRelay/Functions"
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
)

func ViewLogs(GinContext *gin.Context) {

	// Get service 

	Services, Error := Models.ListServices(true) // Will get logs from services in this case

	if Error != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Error retrieving services",
			JSON:    nil,

		})

		return

	}

	// Combine all logs from all services into a single slice

	var AllLogs []Models.Log

	for _, Service := range Services {

		AllLogs = append(AllLogs, Service.Logs...) 

	}

	// Sort logs by created time

	slices.SortFunc(AllLogs, func(A Models.Log, B Models.Log) int {

		if A.Created < B.Created {

			return -1

		} else if A.Created > B.Created {

			return 1
			
		}

		return 0
		
	})

	GinContext.JSON(200, Types.Response{

		Message: fmt.Sprintf("Successfully retrieved %d %s", len(AllLogs), Functions.PluralizeString(len(AllLogs), "log")),
		JSON: AllLogs,

	});

}