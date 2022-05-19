/**
* @Author:  jager
* @Email:   lhj168os@gmail.com
* @File:    log
* @Date:    2021/5/28 1:41 下午
* @package: log
* @Version: v1.0.0
*
* @Description:
*
 */

package logx_bak

//
//import (
//	"fmt"
//	"github.com/zeromicro/go-zero/core/logx"
//	"io"
//	"strings"
//)
//
//const (
//	DebugLevel = "debug"
//	InfoLevel  = "info"
//	WarnLevel  = "warn"
//	ErrorLevel = "error"
//	PanicLevel = "panic"
//
//	FileTypeLog     = "log"
//	FileTypeRequest = "request_log"
//)
//
//var logger *Log
//
//// Log 默认会使用zap作为日志输出引擎. Log集成了日志切割的功能。默认文件大小1024M，自动压缩
//// 最大有3个文件备份，备份保存时间7天。默认不会打印日志被调用的文文件名和位置;
//// 输出:日志默认会被分成五类文件：xxx.log.DEBUG，xxx.log.INFO, xxx.log.WARN, xxx.log.ERROR, xxx.log.request
//// error,panic都会打印在xxx.log.ERROR. 所有的请求都会打在xxx.log.request
//// Adapter:经过比对现在流行的日志库：zap, logrus, zerolog; logrus 虽说格式化，插件化良好，但是
//// 其内部实现锁竞争太过剧烈，性能不好. zap 性能好，格式一般， zerolog性能没有zap好， 相比
//// 来说就没啥突出优点了
//
//type Log struct {
//	Path           string
//	Level          string
//	NeedRequestLog bool
//	adapters       map[string]*zapAdapter
//}
//
//func WriteIO() io.Writer {
//	return logger
//}
//
//func Logger() *Log {
//	return logger
//}
//
//type Option interface {
//	apply(*Log)
//}
//
//type logOptionFunc func(*Log)
//
//func (f logOptionFunc) apply(log *Log) {
//	f(log)
//}
//
//func SetMaxFileSize(size int) Option {
//	return logOptionFunc(func(log *Log) {
//		for i, _ := range log.adapters {
//			log.adapters[i].setMaxFileSize(size)
//		}
//	})
//}
//
//func SetMaxBackups(n int) Option {
//	return logOptionFunc(func(log *Log) {
//		for i, _ := range log.adapters {
//			log.adapters[i].setMaxBackups(n)
//		}
//	})
//}
//
//func SetMaxAge(age int) Option {
//	return logOptionFunc(func(log *Log) {
//		for i, _ := range log.adapters {
//			log.adapters[i].setMaxAge(age)
//		}
//	})
//}
//
//func SetCompress(compress bool) Option {
//	return logOptionFunc(func(log *Log) {
//		for i, _ := range log.adapters {
//			log.adapters[i].setCompress(compress)
//		}
//	})
//}
//
//func SetCaller() Option {
//	return logOptionFunc(func(log *Log) {
//		for i, _ := range log.adapters {
//			log.adapters[i].setCaller(true)
//		}
//	})
//}
//
//func SetRequest() Option {
//	return logOptionFunc(func(log *Log) {
//		log.NeedRequestLog = true
//	})
//}
//
//func DisableZeroLog() {
//	logx.Disable()
//	logx.DisableStat()
//}
//
//func SetupZeroLog(c logx.LogConf) {
//	logx.MustSetup(c)
//}
//
//func SetFileOut(dir string, appName string) Option {
//	var path string
//	if strings.HasSuffix(dir, "/") {
//		path = fmt.Sprintf("%s%s.log", dir, appName)
//	} else {
//		path = fmt.Sprintf("%s/%s.log", dir, appName)
//	}
//
//	return logOptionFunc(func(log *Log) {
//		log.adapters[FileTypeRequest].setFileOut(fmt.Sprintf("%s.request", path))
//		log.adapters[FileTypeLog].setFileOut(path)
//	})
//}
//
//func SetSource(source string) Option {
//	return logOptionFunc(func(log *Log) {
//		for i, _ := range log.adapters {
//			log.adapters[i].setSource(source)
//		}
//	})
//}
//
//// Init init logger
//func Init(level string, options ...Option) {
//	logger = &Log{Level: level}
//	logger.createFiles(level, options...)
//}
//
//func init() {
//	Init(InfoLevel)
//}
//
//// Sync flushes buffer, if any
//func Sync() {
//	if logger == nil {
//		return
//	}
//
//	for _, v := range logger.adapters {
//		v.logger.Sync()
//	}
//}
//
//func (l *Log) isCaller(level string) bool {
//	if v, ok := l.adapters[level]; ok {
//		return v.Caller
//	}
//	return false
//}
//
//func (l *Log) maxFileSize(level string) int {
//	if v, ok := l.adapters[level]; ok {
//		return v.MaxFileSize
//	}
//	return 0
//}
//
//func (l *Log) maxBackups(level string) int {
//	if v, ok := l.adapters[level]; ok {
//		return v.MaxBackups
//	}
//	return 0
//}
//
//func (l *Log) maxAge(level string) int {
//	if v, ok := l.adapters[level]; ok {
//		return v.MaxAge
//	}
//	return 0
//}
//
//func (l *Log) createFiles(level string, options ...Option) {
//	adapters := make(map[string]*zapAdapter, 2)
//	adapters[FileTypeLog] = NewZapAdapter(level)
//	adapters[FileTypeRequest] = NewZapAdapter(InfoLevel)
//	l.adapters = adapters
//
//	for _, opt := range options {
//		opt.apply(l)
//	}
//
//	for _, adapter := range adapters {
//		adapter.Build()
//	}
//
//}
//
//// Debug 使用方法：log.Debug("test")
//func Debug(args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if debugAdapter, ok := logger.adapters[FileTypeLog]; ok {
//		debugAdapter.Debug(args...)
//	}
//}
//
//// Debugf 使用方法：log.Debugf("test:%s", err)
//func Debugf(template string, args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if debugAdapter, ok := logger.adapters[FileTypeLog]; ok {
//		debugAdapter.Debugf(template, args...)
//	}
//}
//
//// Debugw 使用方法：log.Debugw("test", "field1", "value1", "field2", "value2")
//func Debugw(msg string, keysAndValues ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if debugAdapter, ok := logger.adapters[FileTypeLog]; ok {
//		debugAdapter.Debugw(msg, keysAndValues...)
//	}
//}
//
//func Info(args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Info(args...)
//	}
//}
//
//func Infof(template string, args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Infof(template, args...)
//	}
//}
//
//func Infow(msg string, keysAndValues ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Infow(msg, keysAndValues...)
//	}
//}
//
//func Warn(args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Warn(args...)
//	}
//}
//
//func Warnf(template string, args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Warnf(template, args...)
//	}
//}
//
//func Warnw(msg string, keysAndValues ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Warnw(msg, keysAndValues...)
//	}
//}
//
//func Error(args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Error(args...)
//	}
//}
//
//func Errorf(template string, args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Errorf(template, args...)
//	}
//}
//
//func Errorw(msg string, keysAndValues ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Errorw(msg, keysAndValues...)
//	}
//}
//
//func Panic(args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Panic(args...)
//	}
//}
//
//func Panicf(template string, args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Panicf(template, args...)
//	}
//}
//
//func Panicw(msg string, keysAndValues ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Panicw(msg, keysAndValues...)
//	}
//}
//
//func Fatal(args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Fatal(args...)
//	}
//}
//
//func Fatalf(template string, args ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Fatalf(template, args...)
//	}
//}
//
//func Fatalw(msg string, keysAndValues ...interface{}) {
//	if logger == nil {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeLog]; ok {
//		adapter.Fatalw(msg, keysAndValues...)
//	}
//}
//
//func RequestLogInfow(msg string, keysAndValues ...interface{}) {
//	if logger == nil || !logger.NeedRequestLog {
//		return
//	}
//
//	if adapter, ok := logger.adapters[FileTypeRequest]; ok {
//		adapter.Infow(msg, keysAndValues...)
//	}
//}
//
//func (l *Log) Write(p []byte) (n int, err error) {
//	msg := string(p)
//	RequestLogInfow(msg)
//	return len(p), nil
//}
