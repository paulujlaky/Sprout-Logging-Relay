package Routes

import (
	"github.com/gin-gonic/gin"
)

func ViewLogs(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

}
