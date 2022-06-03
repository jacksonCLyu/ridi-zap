package zaplog

import (
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

type zapLoggerBuilder struct {
	zLogger *zapLogger
}

// NewBuilder returns a new zapLoggerBuilder.
func NewBudiler() *zapLoggerBuilder {
	return &zapLoggerBuilder{
		zLogger: &zapLogger{},
	}
}

// ToBuilder returns a new zapLoggerBuilder with the same configuration.
func (z *zapLogger) ToBuilder() *zapLoggerBuilder {
	return &zapLoggerBuilder{
		zLogger: z,
	}
}

// Name set the name of the logger.
func (builder *zapLoggerBuilder) Name(name string) *zapLoggerBuilder {
	builder.zLogger.name = name
	return builder
}

// ZapLogger set a zap.Logger with the given zap.Logger.
func (builder *zapLoggerBuilder) ZapLogger(zapLogger *zap.Logger) *zapLoggerBuilder {
	builder.zLogger.logger = zapLogger
	return builder
}

// LumberLogger set a lumberjack.Logger with the given lumberjack.Logger.
func (builder *zapLoggerBuilder) LumberLogger(lumberLogger *lumberjack.Logger) *zapLoggerBuilder {
	builder.zLogger.lumberLogger = lumberLogger
	return builder
}

// Rotate set the rotate options.
// If rotate is true, the zapLogger will start a goroutine to rotate the lumberjack.Logger in cycle.
// And The Lumberjack.logger will be verified during build time, if the lumberjack.Logger is not valid, Rotate will panic.
func (builder *zapLoggerBuilder) Rotate(rotate bool) *zapLoggerBuilder {
	builder.zLogger.rotate = rotate
	return builder
}

// Delay set the delay options.
func (builder *zapLoggerBuilder) Delay(delay time.Duration) *zapLoggerBuilder {
	builder.zLogger.delay = delay
	return builder
}

// Cycle set the cycle options.
func (builder *zapLoggerBuilder) Cycle(cycle time.Duration) *zapLoggerBuilder {
	builder.zLogger.cycle = cycle
	return builder
}

// Build returns a zapLogger with the given configuration.
func (builder *zapLoggerBuilder) Build() *zapLogger {
	zLogger := builder.zLogger
	if zLogger.rotate {
		if zLogger.lumberLogger == nil {
			panic("Rotate is set, but the Lumberjack Logger configuration is not found")
		}
		if zLogger.delayChan == nil {
			zLogger.delayChan = make(chan struct{})
		}
		// 默认按照小时分割
		if zLogger.delay == -1 {
			currentTime := time.Now().Local()
			if delay, err := GetInitialDelay(currentTime.Hour()+1, 0, 0); err == nil {
				zLogger.delay = delay
			} else {
				zLogger.delay = 0
			}
		}
		if zLogger.cycle == -1 {
			zLogger.cycle = time.Hour
		}
		// 启动日志切割
		go zLogger.startRotateCycling()
		go zLogger.delayDeliver()
	}
	return zLogger
}
