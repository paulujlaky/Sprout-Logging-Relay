package Assets

import (
	"elucid503/SproutLoggingRelay/Middleware"

	// Route Imports from ./Routes

	GlobalRoutes "elucid503/SproutLoggingRelay/Routes"
	APIRoutes "elucid503/SproutLoggingRelay/Routes/API"
	MiscRoutes "elucid503/SproutLoggingRelay/Routes/API/Misc"

	"github.com/gin-gonic/gin"
)

type RouteInfo struct {
	Method  string
	Handler func(*gin.Context)

	// Middleware Opts

	RateLimitConfig *Middleware.RouteRateLimitConfig
	Authorized      bool
}

func GetRoutes() map[string]RouteInfo {

	return map[string]RouteInfo{

		"/": {

			Method:  "GET",
			Handler: GlobalRoutes.Dash,

			RateLimitConfig: nil,
		},

		"/API/": {

			Method:  "GET",
			Handler: APIRoutes.Index,

			RateLimitConfig: nil,
		},

		"/API/Misc/Auth": {

			Method:  "POST",
			Handler: MiscRoutes.Authorize,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,
			},
		},

		"/API/Misc/LogOut": {

			Method:  "POST",
			Handler: MiscRoutes.LogOut,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,
			},
		},
	}

}
