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
	"io"
	"os"
)

type loggerWrite struct {
	file   *os.File
	stdout io.Writer
}

func newLoggerWrite(path string, stdout bool) (*loggerWrite, error) {
	lw := &loggerWrite{}
	if path != "" {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, err
		}
		lw.file = f
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

func (l *loggerWrite) close() (err error) {
	if l.file != nil {
		err = l.file.Close()
	}
	return
}
