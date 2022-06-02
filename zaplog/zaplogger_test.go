package zaplog

import (
	"testing"
	"time"
)

func TestGetLogger(t *testing.T) {
	logger := GetLogger("test")
	logger.Info("test")
}

func TestGetLoggerWithOpts(t *testing.T) {
	logger := GetLoggerWithOptions("test")
	logger.Info("test refPath")
}

func TestGetInitialDelay(t *testing.T) {
	currentTime := time.Now()
	if delay, err := GetInitialDelay(currentTime.Hour()+1, 0, 0); err != nil {
		t.Error(err)
	} else {
		t.Log(delay)
	}
	if delay, err := GetInitialDelay(currentTime.Hour(), currentTime.Minute()+59, currentTime.Second()+59); err != nil {
		t.Error(err)
	} else {
		t.Log(delay)
	}
	if delay, err := GetInitialDelay(0, 0, 0); err != nil {
		t.Error(err)
	} else {
		t.Log(delay)
	}
}
