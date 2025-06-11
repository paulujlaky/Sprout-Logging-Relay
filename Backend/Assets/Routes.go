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

	Path string

	Method  string
	Handler func(*gin.Context)

	// Middleware Opts

	RateLimitConfig *Middleware.RouteRateLimitConfig
	Authorized      bool

}

func GetRoutes() []RouteInfo {

	return []RouteInfo{

		{

			Path:   "/",
			Method:  "GET",

			Handler: GlobalRoutes.Dash,

			RateLimitConfig: nil,

		},

		{

			Path:   "/API/",
			Method:  "GET",

			Handler: APIRoutes.Index,
			RateLimitConfig: nil,

		},

		{
			Path:   "/API/Auth/LogIn",
			Method:  "POST",

			Handler: AuthRoutes.LogIn,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

		},

		{
			Path:   "/API/Auth/LogOut",
			Method:  "POST",

			Handler: AuthRoutes.LogOut,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

		},

		// General Log Routes 

		{

			Path:   "/API/Logs/Watch",
			Method:  "GET",

			Handler: ServiceRoutes.WatchLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},
			Authorized: true,

		},

		// Log-Specific Routes 

		{

			Path:   "/API/Logs/:LogUID",
			Method:  "GET",

			Handler: LogRoutes.FetchLog,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		{

			Path:   "/API/Logs/:LogUID",
			Method:  "DELETE",

			Handler: LogRoutes.DeleteLog,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

			Authorized: true,

		},

		// Service-Specific Log Routes

		{

			Path:   "/API/Services/:ServiceUID/Logs/Create",
			Method:  "POST",

			Handler: LogRoutes.CreateLog,

			Authorized: false, // services must be able to call this route programmatically

		},

		{

			Path:   "/API/Services/:ServiceUID/Logs/View",
			Method:  "GET",

			Handler: ServiceRoutes.ViewServiceSpecificLogs,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 120,
				TimeWindow:         60,

			},

			Authorized: false, // services must be able to call this route programmatically

		},

		// Services

		{

			Path:   "/API/Services/Create",
			Method:  "POST",

			Handler: ServiceRoutes.CreateService,

			RateLimitConfig: &Middleware.RouteRateLimitConfig{

				MaxRequestsAllowed: 90,
				TimeWindow:         60,

			},

			Authorized: true,

		},
		{


			Path:   "/API/Services/List",
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