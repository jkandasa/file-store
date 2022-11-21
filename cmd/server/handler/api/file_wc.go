package api

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	handlerUtils "github.com/jkandasa/file-store/cmd/server/handler/utils"
	"github.com/jkandasa/file-store/pkg/store"
	"github.com/jkandasa/file-store/pkg/types"
	"go.uber.org/zap"
)

// registers wc, freq-words api
func RegisterWordCountRoutes(router *mux.Router) {
	router.HandleFunc("/api/wc/count", wcList).Methods(http.MethodGet)
	router.HandleFunc("/api/wc/freq-words", freqWords).Methods(http.MethodPost)
}

func wcList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// update store
	store.UpdateFilesStore()

	files := store.ListFiles()

	wordCount := uint64(0)
	for _, file := range files {
		// include only text files
		if !strings.HasSuffix(strings.ToLower(file.Name), types.TEXT_FILE_EXTENSION) {
			continue
		}

		filename := filepath.Join(types.STORE_DATA_PATH, file.Name)

		fh, err := os.Open(filename)
		if err != nil {
			zap.L().Error("error on opening a file", zap.String("filename", filename), zap.Error(err))
			handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reader := bufio.NewReader(fh)
		for {
			line, _ := reader.ReadString('\n')
			fields := strings.Fields(line)
			wordCount += uint64(len(fields))
			if line == "" {
				break
			}
		}
	}

	handlerUtils.PostSuccessResponse(w, wordCount)
}

func freqWords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// update store
	store.UpdateFilesStore()

	reqInput := &types.FreqWordsRequest{}
	err := handlerUtils.LoadEntity(w, r, reqInput)
	if err != nil {
		handlerUtils.PostErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := store.ListFiles()

	freqWordsResponse := make(map[string]uint64)
	for _, file := range files {
		// include only text files
		if !strings.HasSuffix(strings.ToLower(file.Name), types.TEXT_FILE_EXTENSION) {
			continue
		}

		filename := filepath.Join(types.STORE_DATA_PATH, file.Name)

		fh, err := os.Open(filename)
		if err != nil {
			zap.L().Error("error on opening a file", zap.String("filename", filename), zap.Error(err))
			handlerUtils.PostErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reader := bufio.NewReader(fh)
		for {
			line, _ := reader.ReadString('\n')
			fields := strings.Fields(line)
			for _, field := range fields {
				freqWordsResponse[field]++
			}
			if line == "" {
				break
			}
		}
	}

	// apply limits and order
	keys := make([]string, 0, len(freqWordsResponse))

	for key := range freqWordsResponse {
		keys = append(keys, key)
	}

	// sort the keys
	if reqInput.OrderBy == types.OrderByASC {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	sort.SliceStable(keys, func(i, j int) bool {
		if reqInput.OrderBy == types.OrderByASC {
			return freqWordsResponse[keys[i]] < freqWordsResponse[keys[j]]
		} else {
			return freqWordsResponse[keys[i]] > freqWordsResponse[keys[j]]
		}
	})

	finalMap := make(map[string]uint64)

	limitKeys := keys
	if reqInput.Limit < uint(len(keys)) {
		limitKeys = keys[:reqInput.Limit]
	}

	for _, k := range limitKeys {
		finalMap[k] = freqWordsResponse[k]
	}

	handlerUtils.PostSuccessResponse(w, finalMap)
}
