package zaplog

import "testing"

func TestGetLogger(t *testing.T) {
	logger := GetLogger("test")
	logger.Info("test")
}

func TestGetLoggerWithOpts(t *testing.T) {
	logger := GetLoggerWithOptions("test")
	logger.Info("test refPath")
}
