/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    logs
 * @Date:    2022/3/21 4:02 PM
 * @package: logs
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package logx

import (
	"fmt"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
)

var logger *loggerSt

type Option struct {
	Level   string
	LogPath string
	Stdout  bool
	Caller  bool
	Source  string
}

func defaultOption() *Option {
	return &Option{
		Level:  "debug",
		Stdout: true,
		Caller: true,
	}
}

type loggerSt struct {
	w  *loggerWrite
	lg zerolog.Logger
}

func callerMarshalFunc(file string, line int) string {
	fs := strings.Split(file, "/")
	fsNum := len(fs)
	if fsNum <= 3 {
		return file + ":" + strconv.Itoa(line)
	}
	return fmt.Sprintf("%s/%s/%s:%d", fs[fsNum-3], fs[fsNum-2], fs[fsNum-1], line)
}

func parseLevel(level string) string {
	switch level {
	case "release":
		return "info"
	case "test":
		return "debug"
	default:
		return level
	}
}

func Init(opfs ...func(opt *Option)) error {
	zerolog.CallerMarshalFunc = callerMarshalFunc
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorFieldName = "err"

	opt := defaultOption()

	for _, opf := range opfs {
		opf(opt)
	}

	if opt.LogPath != "" {
		if !strings.HasSuffix(opt.LogPath, ".log") {
			if opt.Source != "" {
				opt.LogPath = fmt.Sprintf("%s/%s.log", opt.LogPath, opt.Source)
			} else {
				opt.LogPath = fmt.Sprintf("%s/server.log", opt.LogPath)
			}

			opt.LogPath = strings.Replace(opt.LogPath, "//", "/", -1)
		}
	}

	opt.Level = parseLevel(opt.Level)

	lv, err := zerolog.ParseLevel(opt.Level)
	if err != nil {
		return err
	}

	lg, err := newLoggerWrite(opt.LogPath, opt.Stdout)
	if err != nil {
		return err
	}

	cfg := zerolog.New(lg).Level(lv).With().Timestamp().Stack()

	if opt.Caller {
		cfg = cfg.Caller()
	}

	log := cfg.Logger()

	if opt.Source != "" {
		log = log.Hook(&sourcehook{source: opt.Source})
	}

	logger = &loggerSt{
		w:  lg,
		lg: log,
	}
	return nil
}

func Sync() {
	logger.w.close()
}

func Infof(format string, v ...interface{}) {
	Info().Msgf(format, v...)
}

func Info() *zerolog.Event {
	return Logger().Info()
}

func Debugf(format string, v ...interface{}) {
	Debug().Msgf(format, v...)
}

func Debug() *zerolog.Event {
	return Logger().Debug()
}

func Warnf(format string, v ...interface{}) {
	Warn().Msgf(format, v...)
}

func Warn() *zerolog.Event {
	return Logger().Warn()
}

func Errorf(format string, v ...interface{}) {
	Error().Msgf(format, v...)
}

func Error() *zerolog.Event {
	return Logger().Error()
}

func Err(err error) *zerolog.Event {
	return Logger().Err(err)
}

func Fatalf(format string, v ...interface{}) {
	Fatal().Msgf(format, v...)
}

func Fatal() *zerolog.Event {
	return Logger().Fatal()
}

func Panicf(format string, v ...interface{}) {
	Panic().Msgf(format, v...)
}

func Panic() *zerolog.Event {
	return Logger().Panic()
}

//func Tracef(format string, v ...interface{}) {
//	Trace().Msgf(format, v...)
//}
//
//func Trace() *zerolog.Event {
//	return Logger().Trace()
//}

func Logger() *zerolog.Logger {
	if logger == nil {
		Init()
	}
	return &logger.lg
}
