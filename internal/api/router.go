package api

import (
	"banduslib/internal/interfaces"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// route is the information for every URI.
type route struct {
	// Name is the name of this route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// routes is the list of the generated route.
type routes []route

type apiRouter struct {
	apiHandler interfaces.ApiHandler
}

func (a *apiRouter) GetEngine() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Origin", "Content-type"},
	}))

	v1 := router.Group("v1")
	for _, route := range a.getRoutes() {
		switch route.Method {
		case http.MethodGet:
			v1.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			v1.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			v1.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			v1.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

func (a *apiRouter) getRoutes() routes {
	return routes{
		{
			"Index",
			http.MethodGet,
			"/",
			a.apiHandler.Index,
		},

		{
			"AssignTunesToSet",
			http.MethodPut,
			"/sets/:setId/tunes",
			a.apiHandler.AssignTunesToSet,
		},

		{
			"CreateSet",
			http.MethodPost,
			"/sets",
			a.apiHandler.CreateSet,
		},

		{
			"CreateTune",
			http.MethodPost,
			"/tunes",
			a.apiHandler.CreateTune,
		},

		{
			"DeleteSet",
			http.MethodDelete,
			"/sets/:setId",
			a.apiHandler.DeleteSet,
		},

		{
			"DeleteTune",
			http.MethodDelete,
			"/tunes/:tuneId",
			a.apiHandler.DeleteTune,
		},

		{
			"GetSet",
			http.MethodGet,
			"/sets/:setId",
			a.apiHandler.GetSet,
		},

		{
			"GetTune",
			http.MethodGet,
			"/tunes/:tuneId",
			a.apiHandler.GetTune,
		},

		{
			"ListSets",
			http.MethodGet,
			"/sets",
			a.apiHandler.ListSets,
		},

		{
			"ListTunes",
			http.MethodGet,
			"/tunes",
			a.apiHandler.ListTunes,
		},

		{
			"UpdateSet",
			http.MethodPut,
			"/sets/:setId",
			a.apiHandler.UpdateSet,
		},

		{
			"UpdateTune",
			http.MethodPut,
			"/tunes/:tuneId",
			a.apiHandler.UpdateTune,
		},
	}

}

func NewApiRouter(
	apiHandler interfaces.ApiHandler,
) interfaces.ApiRouter {
	return &apiRouter{
		apiHandler: apiHandler,
	}
}
