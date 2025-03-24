package Assets

import (
	"elucid503/SproutLoggingRelay/Middleware"

	// Route Imports from ./Routes

	GlobalRoutes "elucid503/SproutLoggingRelay/Routes"
	APIRoutes "elucid503/SproutLoggingRelay/Routes/API"
	AuthRoutes "elucid503/SproutLoggingRelay/Routes/API/Auth"

	ServiceRoutes "elucid503/SproutLoggingRelay/Routes/API/Services"
	LogRoutes "elucid503/SproutLoggingRelay/Routes/API/Services/Logs"

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

		"/API/Auth/LogIn": {

			Method:  "POST",
			Handler: AuthRoutes.LogIn,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,
			},
		},

		"/API/Auth/LogOut": {

			Method:  "POST",
			Handler: AuthRoutes.LogOut,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,
			},
		},

		// General Log Routes 

		"/API/Services/:ServiceUID/Logs/Create": {

			Method:  "POST",
			Handler: LogRoutes.CreateLog,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,

			},

			Authorized: false, // services must be able to call this route programmatically

		},

		"/API/Services/:ServiceUID/Logs/View": {

			Method:  "GET",
			Handler: LogRoutes.ViewLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 180,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		"/API/Services/:ServiceUID/Logs/Watch": {

			Method:  "GET", // WebSocket Route
			Handler: LogRoutes.WatchLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		// Services
		
		"/API/Services/Create": {

			Method:  "POST",
			Handler: ServiceRoutes.CreateService,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		"/API/Services/List": {

			Method:  "GET",
			Handler: ServiceRoutes.GetServiceList,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 180,
				TimeWindow:         60,

			},

			Authorized: true,

		},

	}

}
