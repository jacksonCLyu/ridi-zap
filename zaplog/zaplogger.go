package zaplog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	LogMaxSize int = 100
	// LogMaxAge max save days of log files
	LogMaxAge int = 365
	// LogMaxBackups max backups of log files
	LogMaxBackups int = 500

	// DefaultLogAggregateDir default directory of log categories
	DefaultLogAggregateDir string = "logs"
	// DefaultLogFileSuffix default suffix string of log file
	DefaultLogFileSuffix string = "log"
	// localTimeLayout default time layout of log file
	localTimeLayout = "2006-01-02 15:04:05"
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
	name         string
	logger       *zap.Logger
	lumberLogger *lumberjack.Logger
	rotate       bool
	delay        time.Duration
	delayChan    chan struct{}
	cycle        time.Duration
}

// ZapLogger returns a new zap logger.
func ZapLogger(name string, opts ...Option) logger.Logger {
	options := &options{
		level:                  zap.NewAtomicLevel(),
		refPath:                name,
		category:               name,
		caller:                 true,
		callerSkip:             1,
		stackTraceLevel:        zap.NewAtomicLevelAt(zap.ErrorLevel),
		isLocalTime:            true,
		isCompress:             true,
		isSampling:             false,
		logRotate:              true,
		logRotateInitialDelay:  -1,
		logRotateCycleDuration: -1,
		fileName:               "",
		maxSize:                LogMaxSize,
		maxAge:                 LogMaxAge,
		maxBackups:             LogMaxBackups,
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

	syncers := make([]zapcore.WriteSyncer, 0, 2)
	// 默认输出到控制台
	syncers = append(syncers, zapcore.AddSync(os.Stdout))
	// 滚动日志
	var lumberLogger *lumberjack.Logger
	if options.logRotate {
		// 日志文件配置
		if options.fileName == "" {
			options.fileName = filepath.Join(env.AppRootPath(), DefaultLogAggregateDir, options.refPath, strings.Join([]string{options.category, DefaultLogFileSuffix}, "."))
		}
		if !filepath.IsAbs(options.fileName) {
			if strings.Contains(options.fileName, DefaultLogAggregateDir) {
				options.fileName = filepath.Join(env.AppRootPath(), options.fileName)
			} else {
				options.fileName = filepath.Join(env.AppRootPath(), DefaultLogAggregateDir, options.fileName)
			}
		}
		// 如果开启切割大小小于默认值，则重置为默认值
		if options.maxSize < LogMaxSize {
			options.maxSize = LogMaxSize
		}
		// 日志切割文档 hook
		lumberLogger = &lumberjack.Logger{
			Filename:   options.fileName,    // 日志文件路径
			MaxSize:    options.maxSize,     // 每个日志文件保存的最大尺寸 单位：M
			MaxAge:     options.maxAge,      // 文件最多保存天数
			MaxBackups: options.maxBackups,  // 日志文件最多保存备份个数
			LocalTime:  options.isLocalTime, // 是否本地时间
			Compress:   options.isCompress,  // 是否压缩
		}
		syncers = append(syncers, zapcore.AddSync(lumberLogger))
	}

	core := zapcore.NewCore(
		encoder,                                 // 输出编码器
		zapcore.NewMultiWriteSyncer(syncers...), // 写入控制台和文件
		options.level,                           // 允许输出的日志级别
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

	if options.stackTraceLevel.Enabled(options.level.Level()) {
		stackTrace := zap.AddStacktrace(options.stackTraceLevel.Level())
		zOpts = append(zOpts, stackTrace)
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
	return NewBudiler().
		Name(name).
		ZapLogger(l).
		Rotate(options.logRotate).
		LumberLogger(lumberLogger).
		Delay(options.logRotateInitialDelay).
		Cycle(options.logRotateCycleDuration).Build()
}

func (z *zapLogger) startRotateCycling() {
	go func() {
		<-z.delayChan
		// 初始化延时结束后，立即执行一次
		if err := z.lumberLogger.Rotate(); err != nil {
			z.Errorf("log rotate logRotateCycleDuration error: %v\n", err)
		}
		// 周期执行
		settingTrickChan := time.Tick(z.cycle)
		for {
			<-settingTrickChan
			if err := z.lumberLogger.Rotate(); err != nil {
				z.Errorf("log rotate logRotateCycleDuration error: %v\n", err)
			}
		}
	}()
	if z.delay > 0 {
		go func() {
			time.Sleep(z.delay)
			close(z.delayChan)
		}()
	} else {
		close(z.delayChan)
	}
}

func getFormatTimeValue(localTime time.Time, hour int, minute int, second int) string {
	minuteStr := strconv.Itoa(minute)
	secondStr := strconv.Itoa(second)
	lMonth := localTime.Month()
	lDay := localTime.Day()
	if minute < 10 {
		minuteStr = "0" + minuteStr
	}
	if second < 10 {
		secondStr = "0" + secondStr
	}
	lMonthStr := strconv.Itoa(int(lMonth))
	lDayStr := strconv.Itoa(lDay)
	if lMonth < 10 {
		lMonthStr = "0" + lMonthStr
	}
	if lDay < 10 {
		lDayStr = "0" + lDayStr
	}
	return fmt.Sprintf("%d-%s-%s %d:%s:%s", localTime.Year(), lMonthStr, lDayStr, hour, minuteStr, secondStr)
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
