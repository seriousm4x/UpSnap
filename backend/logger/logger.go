package logger

import (
	"log"
	"os"
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

func init() {
	// TODO: decide if you want to log to a file to max
	// output := io.MultiWriter(os.Stdout, logFile)
	output := os.Stdout
	Info.SetOutput(output)
	Error.SetOutput(output)
	Debug.SetOutput(output)

	log.SetOutput(Debug.Writer())
	log.SetPrefix("[DEBUG]")
	log.SetFlags(flags)
}
