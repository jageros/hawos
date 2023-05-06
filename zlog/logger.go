/**
* @Author:  jager
* @Email:   lhj168os@gmail.com
* @File:    adapter
* @Date:    2021/5/28 1:41 下午
* @package: log
* @Version: v1.0.0
*
* @Description:
*
 */

package zlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(opfs ...func(opt *Option)) *Logger {
	opt := defaultOption()
	for _, opf := range opfs {
		opf(opt)
	}
	lg := &Logger{}
	lg.build(opt)
	return lg
}

func NewLoggerWithConfig(opt *Option) *Logger {
	lg := &Logger{}
	lg.build(opt)
	return lg
}

// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
func (l *Logger) createLumberjackHook(opt *Option) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   opt.Path,
		MaxSize:    opt.MaxFileSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}
}

func (l *Logger) build(opt *Option) {
	var w zapcore.WriteSyncer

	if opt.StdOut {
		w = zapcore.AddSync(os.Stdout)
	}

	if opt.FileOut {
		if w == nil {
			w = zapcore.AddSync(l.createLumberjackHook(opt))
		} else {
			w = zapcore.NewMultiWriteSyncer(zapcore.AddSync(l.createLumberjackHook(opt)), w)
		}
	}

	var level zapcore.Level
	switch strings.ToLower(opt.Level) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.DebugLevel
	}

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	cnf := zapcore.NewJSONEncoder(conf)
	core := zapcore.NewCore(cnf, w, level)

	l.Logger = zap.New(core)
	if opt.Caller {
		l.Logger = l.WithOptions(zap.AddCaller(), zap.AddCallerSkip(0))
	}

	if opt.Source != "" {
		l.Logger = l.Logger.With(zap.Field{
			Key:    "source",
			Type:   zapcore.StringType,
			String: opt.Source,
		})
	}
}

// ============================ 简化链式调用 =============================

func (l *Logger) With(key string, value interface{}) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Any(key, value))}
}

func (l *Logger) Err(err error) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Error(err))}
}

func (l *Logger) Arg(value interface{}) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Any("arg", value))}
}

func (l *Logger) Namespace(key string) *Logger {
	return &Logger{Logger: l.Logger.With(zap.Namespace(key))}
}

// =========================== gorm log ===============================

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logger.WithOptions(zap.WithCaller(false)).Info(fmt.Sprintf(format, args...))
}
