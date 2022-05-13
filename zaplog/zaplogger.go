package zaplog

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jacksonCLyu/ridi-faces/pkg/env"
	"github.com/jacksonCLyu/ridi-faces/pkg/logger"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// MessageKey zap format message key
	MessageKey string = "msg"
	// TimeKey zap format time key
	TimeKey string = "time"
	// LevelKey zap format level key
	LevelKey string = "level"
	// NameKey zap format name key
	NameKey string = "logger"
	// CallerKey zap format caller key
	CallerKey string = "caller"
	// StacktraceKey zap format stacktrace key
	StacktraceKey string = "stacktrace"

	// LogFormatJSON log output with json type
	LogFormatJSON string = "json"

	// LogMaxSize max size of single log file
	LogMaxSize int = 10
	// LogMaxAge max save days of log files
	LogMaxAge int = 28
	// LogMaxBackups max backups of log files
	LogMaxBackups int = 5

	// DefaultLogAggregateDir default directory of log categories
	DefaultLogAggregateDir string = "logs"
	// DefaultLogFileSuffix default suffix string of log file
	DefaultLogFileSuffix string = "log"
)

var _ logger.Logger = (*zapLogger)(nil)

func (z *zapLogger) Trace(args ...any) {
	z.logger.Sugar().Debug(args...)
}

func (z *zapLogger) Tracef(format string, args ...any) {
	z.logger.Sugar().Debugf(format, args...)
}

func (z *zapLogger) Debug(args ...any) {
	z.logger.Sugar().Debug(args...)
}

func (z *zapLogger) Debugf(format string, args ...any) {
	z.logger.Sugar().Debugf(format, args...)
}

func (z *zapLogger) Info(args ...any) {
	z.logger.Sugar().Info(args...)
}

func (z *zapLogger) Infof(format string, args ...any) {
	z.logger.Sugar().Infof(format, args...)
}

func (z *zapLogger) Warn(args ...any) {
	z.logger.Sugar().Warn(args...)
}

func (z *zapLogger) Warnf(format string, args ...any) {
	z.logger.Sugar().Warnf(format, args...)
}

func (z *zapLogger) Error(args ...any) {
	z.logger.Sugar().Error(args...)
}

func (z *zapLogger) Errorf(format string, args ...any) {
	z.logger.Sugar().Errorf(format, args...)
}

func (z *zapLogger) Fatal(args ...any) {
	z.logger.Sugar().Fatal(args...)
}

func (z *zapLogger) Fatalf(format string, args ...any) {
	z.logger.Sugar().Fatalf(format, args...)
}

type zapLogger struct {
	logger *zap.Logger
}

// ZapLogger returns a new zap logger.
func ZapLogger(opts ...Option) logger.Logger {
	options := &options{
		level:       zap.NewAtomicLevel(),
		refPath:     "app",
		category:    "app",
		caller:      true,
		callerSkip:  2,
		isLocalTime: true,
		isCompress:  true,
		isSampling:  true,
	}
	for _, opt := range opts {
		opt.apply(options)
	}

	// encoder
	var encoder zapcore.Encoder
	if env.IsLocal() {
		//encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		encoder = zapcore.NewConsoleEncoder(NewCustomStdoutEncoderConfig())
	} else {
		//encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		encoder = zapcore.NewConsoleEncoder(NewCustomProductionEncoderConfig())
	}

	// 日志文件配置
	fileName := filepath.Join(env.AppRootPath(), DefaultLogAggregateDir, options.refPath, strings.Join([]string{options.category, DefaultLogFileSuffix}, "."))
	// 日志切割文档 hook
	lumberLogWriter := lumberjack.Logger{
		Filename:   fileName,            // 日志文件路径
		MaxSize:    LogMaxSize,          // 每个日志文件保存的最大尺寸 单位：M
		MaxAge:     LogMaxAge,           // 文件最多保存天数
		MaxBackups: LogMaxBackups,       // 日志文件最多保存备份个数
		LocalTime:  options.isLocalTime, // 是否本地时间
		Compress:   options.isCompress,  // 是否压缩
	}

	core := zapcore.NewCore(
		encoder, // 输出编码器
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&lumberLogWriter)), // 写入控制台和文件
		options.level, // 允许输出的日志级别
	)

	// 构造选项
	zOpts := make([]zap.Option, 0, 3)

	// 开启行号
	if options.caller {
		caller := zap.AddCaller()
		zOpts = append(zOpts, caller)
		callerSkip := zap.AddCallerSkip(options.callerSkip)
		zOpts = append(zOpts, callerSkip)
	}

	// 日志采样
	if options.isSampling {
		samplingConfig := &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
		wrapCore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				samplingConfig.Initial,
				samplingConfig.Thereafter,
				samplerOpts...,
			)
		})
		zOpts = append(zOpts, wrapCore)
	}

	// 构造日志对象
	l := zap.New(core, zOpts...)
	return &zapLogger{logger: l}
}

// NewCustomStdoutEncoderConfig return a custom zapcore encoder config
func NewCustomStdoutEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     MessageKey,
		LevelKey:       LevelKey,
		TimeKey:        TimeKey,
		NameKey:        NameKey,
		CallerKey:      CallerKey,
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,   // second duration encoder
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 短路径编码器(相对路径 + 行号)
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// NewCustomProductionEncoderConfig return a custom zapcore encoder config
func NewCustomProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     MessageKey,
		LevelKey:       LevelKey,
		TimeKey:        TimeKey,
		NameKey:        NameKey,
		CallerKey:      CallerKey,
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,   // second duration encoder
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 短路径编码器(相对路径 + 行号)
		EncodeName:     zapcore.FullNameEncoder,
	}
}
