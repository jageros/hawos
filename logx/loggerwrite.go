/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    loggerwrite
 * @Date:    2022/3/21 6:02 PM
 * @package: logs
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package logx

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type loggerWrite struct {
	file   io.Writer
	stdout io.Writer
}

// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
func createLumberjackHook(path string, maxFileSize, maxBackups, maxAge int, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxFileSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}

func newLoggerWrite(stdout bool, path string, maxFileSize, maxBackups, maxAge int, compress bool) (*loggerWrite, error) {
	lw := &loggerWrite{}
	if path != "" {
		lg := createLumberjackHook(path, maxFileSize, maxBackups, maxAge, compress)
		lw.file = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = lg
			w.TimeFormat = "[01-02 15:04:05]"
		})
	}
	if stdout {
		lw.stdout = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "[01-02 15:04:05]"
		})
	}

	return lw, nil
}

func (l *loggerWrite) Write(data []byte) (n int, err error) {
	if l.stdout != nil {
		n, err = l.stdout.Write(data)
		if err != nil {
			return 0, err
		}
	}
	if l.file != nil {
		n, err = l.file.Write(data)
	}
	return
}

//func (l *loggerWrite) close() (err error) {
//	if l.file != nil {
//		err = l.file.Close()
//	}
//	return
//}
