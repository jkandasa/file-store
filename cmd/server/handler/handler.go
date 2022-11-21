package handler

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	handlerAPI "github.com/jkandasa/file-store/cmd/server/handler/api"
	"github.com/rs/cors"
)

// GetHandler for http access
func GetHandler() (http.Handler, error) {
	router := mux.NewRouter()

	// other routes
	handlerAPI.RegisterVersionRoutes(router)
	handlerAPI.RegisterListFilesRoutes(router)
	handlerAPI.RegisterAddFilesRoutes(router)
	handlerAPI.RegisterRemoveFilesRoutes(router)
	handlerAPI.RegisterWordCountRoutes(router)

	// pre flight middleware
	withCors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	withPreflight := withCors.Handler(router)

	// include gzip middleware
	withGzip := gziphandler.GzipHandler(withPreflight)

	return withGzip, nil
}
