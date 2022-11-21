package api

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	handlerUtils "github.com/jkandasa/file-store/cmd/server/handler/utils"
	"github.com/jkandasa/file-store/pkg/store"
	"github.com/jkandasa/file-store/pkg/types"
	"github.com/jkandasa/file-store/pkg/utils"
)

// registers remove api
func RegisterRemoveFilesRoutes(router *mux.Router) {
	router.HandleFunc("/api/file/remove", removeFiles).Methods(http.MethodPost)
}

func removeFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// update store
	defer store.UpdateFilesStore()

	files := []string{}
	err := handlerUtils.LoadEntity(w, r, &files)
	if err != nil {
		handlerUtils.PostErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// remove files
	errs := make([]error, 0)
	for _, filename := range files {
		_filename := filepath.Join(types.HOME_PATH, filename)
		if utils.IsFileExists(_filename) {
			err := utils.RemoveFile(_filename)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// TODO: send all the errors
	if len(errs) > 0 {
		handlerUtils.PostErrorResponse(w, errs[0].Error(), http.StatusInternalServerError)
		return
	}

	handlerUtils.PostSuccessResponse(w, "files removed successfully")
}
