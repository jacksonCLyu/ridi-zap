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
	isLocalTime bool
	isCompress  bool
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
