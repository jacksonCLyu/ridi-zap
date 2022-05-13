package zaplog

import "go.uber.org/zap"

// Option zap logger option interface
type Option interface {
	apply(zapOpts *options)
}

type options struct {
	level       zap.AtomicLevel
	refPath     string
	category    string
	caller      bool
	callerSkip  int
	isLocalTime bool
	isCompress  bool
	isSampling  bool
}

// ApplyFunc zap logger options apply func
type ApplyFunc func(zapOpts *options)

func (f ApplyFunc) apply(zapOpts *options) {
	f(zapOpts)
}

// WithLevel return logger with the given level
func WithLevel(level zap.AtomicLevel) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.level = level
	})
}

// WithRefPath return logger with the given ref path
func WithRefPath(logRefPath string) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.refPath = logRefPath
	})
}

// WithCategory return logger with category
func WithCategory(category string) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.category = category
		if zapOpts.refPath == "" {
			zapOpts.refPath = category
		}
	})
}

// AddCaller return logger with caller
func AddCaller() Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.caller = true
	})
}

// WithCaller return logger with caller option
func WithCaller(caller bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.caller = caller
	})
}

// WithCallerSkip return logger with caller skip
func WithCallerSkip(skip int) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.callerSkip = skip
	})
}

// LocalDateTime return logger with localDateTime
func LocalDateTime(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.isLocalTime = is
	})
}

// Compress return logger if compress needed
func Compress(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.isCompress = is
	})
}

// Sampling return logger if sampling needed
func Sampling(is bool) Option {
	return ApplyFunc(func(zapOpts *options) {
		zapOpts.isSampling = is
	})
}
