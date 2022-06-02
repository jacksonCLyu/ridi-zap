package zaplog

import (
	"go.uber.org/zap"
	"time"
)

// Option zap logger option interface
type Option interface {
	apply(zapOpts *options)
}

type options struct {
	level                  zap.AtomicLevel
	refPath                string
	category               string
	caller                 bool
	callerSkip             int
	stackTraceLevel        zap.AtomicLevel
	isLocalTime            bool
	isCompress             bool
	isSampling             bool
	logRotate              bool
	logRotateInitialDelay  time.Duration
	logRotateCycleDuration time.Duration
	fileName               string
	maxSize                int
	maxAge                 int
	maxBackups             int
}

// ApplyFunc zap logger options apply func
type ApplyFunc func(zapOpts *options)

func (f ApplyFunc) apply(zapOpts *options) {
	f(zapOpts)
}

// WithLevel return logger with the given level default is zap.NewAtomicLevel()
func WithLevel(level zap.AtomicLevel) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.level = level
	})
}

// WithRefPath return logger with the given ref path default is "app"
func WithRefPath(logRefPath string) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.refPath = logRefPath
	})
}

// WithCategory return logger with category default is "app"
func WithCategory(category string) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.category = category
		if zapOpts.refPath == "" {
			zapOpts.refPath = category
		}
	})
}

// AddCaller return logger with caller
// caller is a bool option that determines whether the line number of the caller
// is logged.
func AddCaller() Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.caller = true
	})
}

// WithCaller return logger with caller option
// caller is a bool option that determines whether the line number of the caller
// is logged. The default is true.
func WithCaller(caller bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.caller = caller
	})
}

// WithCallerSkip return logger with caller skip
// callerSkip is the number of stack frames to ascend, starting from zaplogger.go.
// The default is 1.
func WithCallerSkip(skip int) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.callerSkip = skip
	})
}

// WithStackTraceLevel return logger with stackTraceLevel
// stackTraceLevel is used to determine whether to include stack trace in the log output.
// The default is ErrorLevel.
func WithStackTraceLevel(level zap.AtomicLevel) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.stackTraceLevel = level
	})
}

// LocalDateTime return logger with localDateTime default is true
func LocalDateTime(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.isLocalTime = is
	})
}

// Compress return logger if compress needed default is true
func Compress(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.isCompress = is
	})
}

// Sampling return logger if sampling needed default is false
func Sampling(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.isSampling = is
	})
}

// LogRotate return logger if logRotate needed default is true
func LogRotate(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.logRotate = is
	})
}

// LogRotateInitialDelay return logger with initial logRotateInitialDelay time.Duration
// if logRotate is true, initialDelay is the time to wait before the first log file rotation.
// The default is the Duration between time.Now() and the parse time of the time.Now().Hour():59:59.
func LogRotateInitialDelay(delay time.Duration) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.logRotateInitialDelay = delay
	})
}

// LogRotateCycle return logger with logRotateCycleDuration
func LogRotateCycle(cycle time.Duration) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.logRotateCycleDuration = cycle
	})
}

// FileName return logger with logFileName
func FileName(fileName string) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.fileName = fileName
	})
}

// MaxSize return logger with maxSize
// if logRotate is true, maxSize is the maximum size of the log file in megabytes.
// The default is 100 megabytes.
func MaxSize(maxSize int) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.maxSize = maxSize
	})
}

// MaxAge return logger with maxAge
func MaxAge(maxAge int) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.maxAge = maxAge
	})
}

// MaxBackups return logger with maxBackups
func MaxBackups(maxBackups int) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.maxBackups = maxBackups
	})
}
