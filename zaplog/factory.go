package zaplog

import (
	"sync"

	"github.com/jacksonCLyu/ridi-faces/pkg/logger"
)

var logMemCache sync.Map

// GetLogger returns a logger with the given name.
func GetLogger(name string) logger.Logger {
	l, _ := logMemCache.LoadOrStore(name, ZapLogger(WithCategory(name), WithRefPath(name)))
	return l.(logger.Logger)
}

// GetLoggerWithOptions returns a logger with the given name and options.
func GetLoggerWithOptions(name string, opts ...Option) logger.Logger {
	if len(opts) == 0 {
		opts = append(opts, WithCategory(name), WithRefPath(name))
	}
	l, _ := logMemCache.LoadOrStore(name, ZapLogger(opts...))
	return l.(logger.Logger)
}
