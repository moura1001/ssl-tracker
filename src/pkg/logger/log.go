package logger

import (
	"io"
	"os"

	kitlog "github.com/go-kit/log"
)

var logger kitlog.Logger

func Init() {
	var (
		logout io.Writer
	)
	logout = os.Stdout

	w := kitlog.NewSyncWriter(logout)
	logger = kitlog.NewLogfmtLogger(w)
}

func Log(keyvals ...interface{}) {
	logger.Log(keyvals...)
}
