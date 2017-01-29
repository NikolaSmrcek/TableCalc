package routes

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v5"
)

var RedisClient *redis.Client;

// Route struct taking name of the route, method (http method), path and handler function.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
// List of Routes.
type Routes []Route

var routes = Routes{
	Route{
		"GetRedisValueByKey",
		"GET",
		"/redis/{key}",
		GetRedisValueByKey,
	},
	Route{
		"SetRedisValueByKey",
		"POST",
		"/setCells",
		SetRedisValueByKey,
	},
	Route{
		"GetRedisData",
		"GET",
		"/getInitData",
		GetRedisData,
	},
	Route{
		"RemoveFromList",
		"POST",
		"/removeFromList",
		RemoveFromList,
	},
	Route{
		"AddToList",
		"POST",
		"/addToList",
		AddToList,
	},
	Route{
		"LoaderFunction",
		"POST",
		"/upload",
		LoaderFunction,
	},
}
// Used for logging function execution time and to know which route has been called.
func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
// Used for serving static files. Mostly for includes inside of the html files.
func serveStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"styles":           staticDirectory + "/styles/",
		"scripts":          staticDirectory + "/scripts/",
	};
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}

// Creating new gorilla mux router that will server routes and static files.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")));
	serveStatic(router, "/public/");
	return router;
}

func InitRouter(rc *redis.Client) {
	RedisClient = rc;
	InitRedisColsAndRows()
	InitReduceMapper()
}