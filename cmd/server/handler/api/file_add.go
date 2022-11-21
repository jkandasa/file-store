package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	handlerUtils "github.com/jkandasa/file-store/cmd/server/handler/utils"
	"github.com/jkandasa/file-store/pkg/store"
	"github.com/jkandasa/file-store/pkg/types"
	"github.com/jkandasa/file-store/pkg/utils"
	"go.uber.org/zap"
)

// registers files add, update api
func RegisterAddFilesRoutes(router *mux.Router) {
	router.HandleFunc("/api/file/add", addFiles).Methods(http.MethodPost)
	router.HandleFunc("/api/file/update", updateFiles).Methods(http.MethodPost)
	router.HandleFunc("/api/file/bucket", addDataBucket).Methods(http.MethodPost)
}

func addFiles(w http.ResponseWriter, r *http.Request) {
	modifyFiles(w, r, false)
}

func updateFiles(w http.ResponseWriter, r *http.Request) {
	modifyFiles(w, r, true)
}

func modifyFiles(w http.ResponseWriter, r *http.Request, supportUpdate bool) {
	w.Header().Set("Content-Type", "application/json")

	// update store
	store.UpdateFilesStore()

	reqFile := &types.File{}
	err := handlerUtils.LoadEntity(w, r, reqFile)
	if err != nil {
		handlerUtils.PostErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// verify the content or file already exists
	file := store.GetByName(reqFile.Name)
	if file != nil {
		if supportUpdate {
			if file.MD5Hash != reqFile.MD5Hash {
				handlerUtils.PostSuccessResponse(w, string(types.FileResponseNotAvailable))
			} else {
				handlerUtils.PostSuccessResponse(w, string(types.FileResponseUpToDate))
			}
			return
		}
		handlerUtils.PostSuccessResponse(w, string(types.FileResponseNameExists))
		return
	}

	// if content of the file already present, clone it
	file = store.GetByHash(reqFile.MD5Hash)
	if file != nil {
		srcFile := filepath.Join(types.STORE_DATA_PATH, file.Name)
		dstFile := filepath.Join(types.STORE_DATA_PATH, reqFile.Name)
		err = utils.CopyFile(srcFile, dstFile, true)
		if err != nil {
			handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handlerUtils.PostSuccessResponse(w, string(types.FileResponseCloned))
		return
	}

	handlerUtils.PostSuccessResponse(w, string(types.FileResponseNotAvailable))
}

func addDataBucket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fileBucket := &types.Bucket{}
	err := handlerUtils.LoadEntity(w, r, fileBucket)
	if err != nil {
		handlerUtils.PostErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add it in to tmp location
	filename := fmt.Sprintf("%s_%s", fileBucket.UUID, fileBucket.File.Name)
	err = utils.AppendFile(filepath.Join(types.STORE_DATA_PATH, types.TMP_PATH), filename, fileBucket.Data, fileBucket.Index*fileBucket.BucketSize)
	if err != nil {
		handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// file copy done, it is time to place it on the actual location
	if fileBucket.IsLastBucket {
		srcFile := filepath.Join(types.STORE_DATA_PATH, types.TMP_PATH, filename)
		dstFile := filepath.Join(types.STORE_DATA_PATH, fileBucket.File.Name)
		err = utils.CopyFile(srcFile, dstFile, true)
		if err != nil {
			zap.L().Error("error on copying the file to actual location", zap.String("filename", fileBucket.File.Name), zap.Error(err))
			handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = utils.RemoveFile(srcFile)
		if err != nil {
			zap.L().Error("error on removing temporary file", zap.String("filename", srcFile), zap.Error(err))
			handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
