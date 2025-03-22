package Routes

import (
	"elucid503/SproutLoggingRelay/Types"

	"github.com/gin-gonic/gin"
)

func LogOut(GinContext *gin.Context) {

	GinContext.SetCookie("Sprout-JWT", "", -1, "/", "", false, true)

	GinContext.JSON(200, Types.Response{

		Message: "Logged out",
	})

}
