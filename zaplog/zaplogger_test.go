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

func TestStartRotateCycling(t *testing.T) {
	currentTime := time.Now()
	h, m, s := currentTime.Clock()
	delay, err := GetInitialDelay(h, m, s+30)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(delay)
	}
	// _ = GetLoggerWithOptions("test", LogRotateInitialDelay(delay), LogRotateCycle(time.Second*30))
	logger := GetLoggerWithOptions("test", LogRotateInitialDelay(delay), LogRotateCycle(time.Second*30))
	tickerChan := time.Tick(2 * time.Second)
	go func() {
		for {
			<-tickerChan
			logger.Info("test")
		}
	}()
	time.Sleep(time.Minute * 2)
}
