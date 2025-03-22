package Routes

import (
	"github.com/gin-gonic/gin"
)

func WatchLogs(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

}
