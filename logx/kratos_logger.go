package logx

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
)

func (l *loggerSt) Log(level log.Level, keyvals ...interface{}) error {
	var lg *zerolog.Event
	switch level {
	case log.LevelInfo:
		lg = l.lg.Info()
	case log.LevelWarn:
		lg = l.lg.Warn()
	case log.LevelError:
		lg = l.lg.Error()
	case log.LevelFatal:
		lg = l.lg.Fatal()
	default:
		lg = l.lg.Debug()
	}

	vNum := len(keyvals)
	for i := 0; i < vNum-1; i += 2 {
		key, _ := keyvals[i].(string)
		val := keyvals[i+1]
		lg = lg.Interface(key, val)
	}

	lg.Msg("kratos-log")

	return nil
}

func KratosLogger() *loggerSt {
	if logger == nil {
		Init()
	}
	return logger
}
