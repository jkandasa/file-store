package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jkandasa/file-store/pkg/types"
	"go.uber.org/zap"
)

func WriteResponse(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		zap.L().Error("error on writing response", zap.Error(err))
		return
	}
}

func PostErrorResponse(w http.ResponseWriter, message string, code int) {
	response := &types.HttpResponse{
		Success: false,
		Message: message,
	}
	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, string(out), code)
}

func PostSuccessResponse(w http.ResponseWriter, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteResponse(w, out)
}

// convert json to struct based entity
func LoadEntity(w http.ResponseWriter, r *http.Request, entity interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	d, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, &entity)
	if err != nil {
		return err
	}
	return nil
}
