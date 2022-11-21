package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jkandasa/file-store/cmd/server/handler"
	httpLogger "github.com/jkandasa/file-store/cmd/server/handler/logger"
	"github.com/jkandasa/file-store/pkg/store"
	"github.com/jkandasa/file-store/pkg/utils"
	"github.com/jkandasa/file-store/pkg/version"
	"go.uber.org/zap"
)

const (
	defaultReadtimeout = time.Second * 60
)

func main() {
	// init logger
	initLogger()

	httpPort := flag.Int("port", 8080, "http port to serve")
	printVersion := flag.Bool("version", false, "prints version and exits")
	flag.Parse()

	if *printVersion {
		fmt.Println(version.Get().String())
		os.Exit(0)
	}

	zap.L().Info("version details", zap.Any("version", version.Get()))
	// update files store
	store.UpdateFilesStore()

	_handler, err := handler.GetHandler()
	if err != nil {
		zap.L().Fatal("error on getting handlers", zap.Error(err))
	}

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", *httpPort)
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
