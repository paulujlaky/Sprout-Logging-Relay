package Routes

import (
	"elucid503/SproutLoggingRelay/Models"

	"github.com/gin-gonic/gin"
)

func CreateService(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	// Check body

	if (Body["Name"] == nil) || (Body["Description"] == nil) || (Body["Icon"] == nil) {

		GinContext.JSON(400, gin.H{

			"Message": "Missing required fields",

		})

		return
		
	}

	// Parse body

	Name := Body["Name"].(string)
	Description := Body["Description"].(string)
	Icon := Body["Icon"].(string)

	// Create service

	Service := Models.NewService(Name, Description, Icon)

	GinContext.JSON(200, gin.H{

		"Message": "Service created",

		"JSON": Service,

	})

}