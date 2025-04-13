package Assets

import (
	"elucid503/SproutLoggingRelay/Middleware"

	// Route Imports from ./Routes

	GlobalRoutes "elucid503/SproutLoggingRelay/Routes"
	APIRoutes "elucid503/SproutLoggingRelay/Routes/API"
	AuthRoutes "elucid503/SproutLoggingRelay/Routes/API/Auth"

	LogRoutes "elucid503/SproutLoggingRelay/Routes/API/Logs"
	ServiceRoutes "elucid503/SproutLoggingRelay/Routes/API/Services"

	"github.com/gin-gonic/gin"
)

type RouteInfo struct {

	Method  string
	Handler func(*gin.Context)

	// Middleware Opts

	RateLimitConfig *Middleware.RouteRateLimitConfig
	Authorized      bool

	// Other 

	UseMutex bool // if true, this route will use a mutex to lock the route while it is being processed

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

		"/API/Logs/Watch": {

			Method:  "GET",
			Handler: ServiceRoutes.WatchLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		"/API/Logs/View": {

			Method:  "GET",
			Handler: LogRoutes.ViewLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 180,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		// Service-Specific Log Routes

		"/API/Services/:ServiceUID/Logs/Create": {

			Method:  "POST",
			Handler: LogRoutes.CreateLog,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,

			},

			Authorized:false, // services must be able to call this route programmatically

			UseMutex: true,

		},

		"/API/Services/:ServiceUID/Logs/View": {

			Method:  "GET",
			Handler: ServiceRoutes.ViewServiceSpecificLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,

			},

			Authorized:false, // services must be able to call this route programmatically

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
