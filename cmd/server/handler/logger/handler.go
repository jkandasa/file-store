package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jkandasa/file-store/pkg/utils"
	"go.uber.org/zap"
)

const callerSkipLevel = int(4)

var debugMessages = []string{"http: TLS handshake error from"}

// myLogger to control http server logs
// implementation taken from https://github.com/uber-go/zap/blob/v1.17.0/global.go#L77
type myLogger struct {
	logger *zap.Logger
}

func GetHttpLogger(level, encoding string, enableStacktrace bool) *myLogger {
	return &myLogger{logger: utils.GetLogger(level, encoding, false, callerSkipLevel, enableStacktrace)}
}

func (ml *myLogger) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	logMsg := ml.fmtMsg(string(p))

	// debug messages
	for _, debugContent := range debugMessages {
		if strings.Contains(logMsg, debugContent) {
			ml.logger.Debug(logMsg)
			return len(logMsg), nil
		}
	}

	// info message
	ml.logger.Info(logMsg)
	return len(logMsg), nil
}

func (ml *myLogger) fmtMsg(msg string) string {
	m := fmt.Sprintf("[HANDLER] %s", msg)
	return strings.TrimSuffix(m, "\n")
}
