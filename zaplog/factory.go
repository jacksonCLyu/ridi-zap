package zaplog

import (
	"github.com/jacksonCLyu/ridi-faces/pkg/logger"
	"sync"
)

var logMemCache = make(map[string]logger.Logger)
var mu sync.RWMutex

// GetLogger returns a logger with the given name.
func GetLogger(name string) logger.Logger {
	if check(name) {
		return logMemCache[name]
	}
	return store(name)
}

// GetLoggerWithOptions returns a logger with the given name and options.
func GetLoggerWithOptions(name string, opts ...Option) logger.Logger {
	if check(name) {
		return logMemCache[name]
	}
	return store(name, opts...)
}

func check(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := logMemCache[name]
	return ok
}

func store(name string, opts ...Option) logger.Logger {
	mu.Lock()
	defer mu.Unlock()
	l, ok := logMemCache[name]
	if ok {
		return l
	}
	l = ZapLogger(name, opts...)
	logMemCache[name] = l
	return l
}
