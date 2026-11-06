package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

var consoleHandler slog.Handler

// ConfigureConsole sets console output for UpSnap logs. PocketBase persistence
// remains controlled by PocketBase's built-in log settings.
func ConfigureConsole(level string) error {
	level = strings.TrimSpace(level)
	if level == "" {
		level = "INFO"
	}

	if strings.EqualFold(level, "OFF") {
		consoleHandler = nil
		return nil
	}

	var parsed slog.Level
	if err := parsed.UnmarshalText([]byte(strings.ToUpper(level))); err != nil {
		return fmt.Errorf("invalid UPSNAP_LOG_LEVEL %q (expected DEBUG, INFO, WARN, ERROR, or OFF)", level)
	}

	consoleHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parsed,
	})
	return nil
}

func relativeSource(file string) string {
	file = strings.ReplaceAll(file, "\\", "/")
	if index := strings.LastIndex(file, "/backend/"); index >= 0 {
		return file[index+1:]
	}
	return file
}

// Logger writes through PocketBase's logger and mirrors UpSnap logs to stdout.
func Logger(app core.App) *slog.Logger {
	handlers := []slog.Handler{app.Logger().Handler()}
	if consoleHandler != nil && !app.IsDev() {
		handlers = append(handlers, consoleHandler)
	}
	return slog.New(multiHandler(handlers)).With("upsnap", true)
}

type multiHandler []slog.Handler

func (h multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h multiHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := frames.Next()
		record.AddAttrs(slog.String(slog.SourceKey, relativeSource(frame.File)+":"+strconv.Itoa(frame.Line)))
	}

	var errs []error
	for _, handler := range h {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record.Clone()); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make(multiHandler, len(h))
	for i, handler := range h {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return handlers
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	handlers := make(multiHandler, len(h))
	for i, handler := range h {
		handlers[i] = handler.WithGroup(name)
	}
	return handlers
}
