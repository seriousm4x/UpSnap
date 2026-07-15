package logger

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	flags = log.Ldate | log.Ltime | log.Lshortfile
	// maximum number of log entries kept in memory for the /api/upsnap/logs endpoint
	bufferSize = 1000
)

// Entry is a single captured log line, exposed to the frontend via the logs API.
type Entry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

// ringBuffer is an in-memory writer that keeps the last bufferSize log entries.
type ringBuffer struct {
	mu      sync.Mutex
	entries []Entry
}

func (r *ringBuffer) add(level string, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, Entry{Time: time.Now(), Level: level, Message: message})
	if len(r.entries) > bufferSize {
		r.entries = r.entries[len(r.entries)-bufferSize:]
	}
}

// levelWriter adapts the ring buffer to io.Writer for a specific log level, so it
// can be plugged into a standard log.Logger via io.MultiWriter.
type levelWriter struct {
	level string
}

func (w levelWriter) Write(p []byte) (int, error) {
	buffer.add(w.level, string(p))
	return len(p), nil
}

var buffer = &ringBuffer{}

var (
	Info    = log.New(os.Stdout, "[INFO] ", flags)
	Debug   = log.New(os.Stdout, "[DEBUG] ", flags)
	Warning = log.New(os.Stdout, "[WARNING] ", flags)
	Error   = log.New(os.Stderr, "[ERROR] ", flags)
)

// Entries returns a copy of the currently buffered log entries, oldest first.
func Entries() []Entry {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	entries := make([]Entry, len(buffer.entries))
	copy(entries, buffer.entries)
	return entries
}

func init() {
	Info.SetOutput(io.MultiWriter(os.Stdout, levelWriter{"INFO"}))
	Error.SetOutput(io.MultiWriter(os.Stderr, levelWriter{"ERROR"}))
	Debug.SetOutput(io.MultiWriter(os.Stdout, levelWriter{"DEBUG"}))
	Warning.SetOutput(io.MultiWriter(os.Stdout, levelWriter{"WARNING"}))

	log.SetOutput(Debug.Writer())
	log.SetPrefix("[DEBUG]")
	log.SetFlags(flags)
}
