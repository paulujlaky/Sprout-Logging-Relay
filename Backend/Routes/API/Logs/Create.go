package Routes

import (
	"github.com/gin-gonic/gin"
)

func CreateLog(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

}
