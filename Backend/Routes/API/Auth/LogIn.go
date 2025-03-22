package Routes

import (
	"elucid503/SproutLoggingRelay/Middleware"
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"slices"

	"github.com/gin-gonic/gin"
)

func LogIn(GinContext *gin.Context) {

	// Get request body

	Body := gin.H{}

	GinContext.BindJSON(&Body)

	// Get token

	Token, TokenExists := Body["Token"].(string)

	if !TokenExists {

		// Try cookie

		Token, _ = GinContext.Cookie("Sprout-JWT")

	}

	if Token == "" {

		GinContext.JSON(401, Types.Response{

			Message: "Unauthorized",
		})

		return

	}

	// Get account

	Account, CachedExists := Middleware.CheckAccountFromCache(Token)

	if !CachedExists {

		var FetchErr error

		Account, FetchErr = Models.GetSproutAccountByToken(Token)

		if FetchErr != nil {

			GinContext.JSON(401, Types.Response{

				Message: "Unauthorized", // No other options
			})

		}

	}

	if !slices.Contains(Account.Flags, "Developer") && !slices.Contains(Account.Flags, "CDNAccess") {

		GinContext.JSON(401, Types.Response{

			Message: "Unauthorized",
			JSON: &gin.H{

				"WasFlagError": true,
			},
		})

		return

	}

	Middleware.AddAccountToCache(Token, Account)

	GinContext.SetCookie("Sprout-JWT", Token, 60*60*24*7, "/", "", false, true)

	GinContext.JSON(200, Types.Response{

		Message: "Authorized",
		JSON:    Account,
	})

}
