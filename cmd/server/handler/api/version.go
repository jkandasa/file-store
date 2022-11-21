package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	handlerUtils "github.com/jkandasa/file-store/cmd/server/handler/utils"
	"github.com/jkandasa/file-store/pkg/version"
)

// registers version api
func RegisterVersionRoutes(router *mux.Router) {
	router.HandleFunc("/api/version", getVersion).Methods(http.MethodGet)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	v := version.Get()
	od, err := json.Marshal(&v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	handlerUtils.WriteResponse(w, od)
}
