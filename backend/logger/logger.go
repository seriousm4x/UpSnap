package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/pocketbase/pocketbase/core"
)

const (
	flags = log.Ldate | log.Ltime | log.Lshortfile
)

var (
	Info    = log.New(os.Stdout, "[INFO] ", flags)
	Debug   = log.New(os.Stdout, "[DEBUG] ", flags)
	Warning = log.New(os.Stdout, "[WARNING] ", flags)
	Error   = log.New(os.Stderr, "[ERROR] ", flags)
)


func InitLogger(app core.App) {

	infoHandler := slog.NewLogLogger(app.Logger().Handler(), slog.LevelInfo)
	errorHandler := slog.NewLogLogger(app.Logger().Handler(), slog.LevelError)
	debugHandler := slog.NewLogLogger(app.Logger().Handler(), slog.LevelDebug)
	warnHandler := slog.NewLogLogger(app.Logger().Handler(), slog.LevelWarn)

	Info.SetOutput(infoHandler.Writer())
	Error.SetOutput(errorHandler.Writer())
	Debug.SetOutput(debugHandler.Writer())
	Warning.SetOutput(warnHandler.Writer())

	log.SetOutput(Debug.Writer())
	log.SetPrefix("[DEBUG]")
	log.SetFlags(flags)
}
