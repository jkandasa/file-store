package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	handlerUtils "github.com/jkandasa/file-store/cmd/server/handler/utils"
	"github.com/jkandasa/file-store/pkg/store"
)

// registers list api
func RegisterListFilesRoutes(router *mux.Router) {
	router.HandleFunc("/api/file/list", listFiles).Methods(http.MethodGet)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// update files cache store
	store.UpdateFilesStore()

	files := store.ListFiles()
	od, err := json.Marshal(&files)
	if err != nil {
		handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	handlerUtils.WriteResponse(w, od)
}
