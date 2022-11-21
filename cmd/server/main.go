package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jkandasa/file-store/cmd/server/handler"
	httpLogger "github.com/jkandasa/file-store/cmd/server/handler/logger"
	"github.com/jkandasa/file-store/pkg/store"
	"github.com/jkandasa/file-store/pkg/utils"
	"go.uber.org/zap"
)

const (
	defaultReadtimeout = time.Second * 60
)

func main() {
	// init logger
	initLogger()

	// update files store
	store.UpdateFilesStore()

	_handler, err := handler.GetHandler()
	if err != nil {
		zap.L().Fatal("error on getting handlers", zap.Error(err))
	}

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	zap.L().Info("listening HTTP service on", zap.String("address", addr))
	server := &http.Server{
		ReadTimeout: defaultReadtimeout,
		Addr:        addr,
		Handler:     _handler,
		ErrorLog:    log.New(httpLogger.GetHttpLogger("debug", "console", false), "", 0),
	}

	err = server.ListenAndServe()
	if err != nil {
		zap.L().Fatal("error on starting http handler", zap.Error(err))
	}
}

func initLogger() {
	logger := utils.GetLogger("debug", "console", false, 0, false)
	zap.ReplaceGlobals(logger)
}
