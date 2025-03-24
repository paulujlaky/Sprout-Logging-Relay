package Routes

import (
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func Index(GinContext *gin.Context) {

	GinContext.JSON(200, Types.Response{

		Message: "Welcome to the Sprout Logging Relay API! Unfortunately, this API is not available for public use. Learn more about Sprout at https://sprout.software. If you are a developer and would like to use this service, hit up developers@sprout.software",
	})

}