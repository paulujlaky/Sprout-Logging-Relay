package Routes

import (
	"github.com/gin-gonic/gin"
)

func Dash(GinContext *gin.Context) {

	// Return the frontend Dash.html file

	GinContext.File("../Frontend/Pages/Dash.html")

}
