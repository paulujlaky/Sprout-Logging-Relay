package Routes

import (
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func Index(GinContext *gin.Context) {

	GinContext.JSON(200, Types.Response{

		Message: "Welcome to the Sprout Logging Relay API!",
	})

}
