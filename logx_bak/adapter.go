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

package logx_bak

//
//import (
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//	"gopkg.in/natefinch/lumberjack.v2"
//	"os"
//)
//
//type zapAdapter struct {
//	Path        string // 文件绝对地址，如：/home/homework/neso/file.log
//	Level       string // 日志输出的级别
//	MaxFileSize int    // 日志文件大小的最大值，单位(M)
//	MaxBackups  int    // 最多保留备份数
//	MaxAge      int    // 日志文件保存的时间，单位(天)
//	Compress    bool   // 是否压缩
//	Caller      bool   // 日志是否需要显示调用位置
//	FileOut     bool   // 是否需要文件输出
//	Source      string // 标志
//
//	logger *zap.Logger
//	sugar  *zap.SugaredLogger
//}
//
//func (z *zapAdapter) setMaxFileSize(size int) {
//	z.MaxFileSize = size
//}
//
//func (z *zapAdapter) setMaxBackups(n int) {
//	z.MaxBackups = n
//}
//
//func (z *zapAdapter) setMaxAge(age int) {
//	z.MaxAge = age
//}
//
//func (z *zapAdapter) setCompress(compress bool) {
//	z.Compress = compress
//}
//
//func (z *zapAdapter) setCaller(caller bool) {
//	z.Caller = caller
//}
//
//func (z *zapAdapter) setFileOut(path string) {
//	z.FileOut = true
//	z.Path = path
//}
//
//func (z *zapAdapter) setSource(source string) {
//	z.Source = source
//}
//
//func NewZapAdapter(level string) *zapAdapter {
//	return &zapAdapter{
//		Level:       level,
//		MaxFileSize: 1024,
//		MaxBackups:  3,
//		MaxAge:      7,
//		Compress:    true,
//		Caller:      false,
//	}
//}
//
//// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
//func (z *zapAdapter) createLumberjackHook() *lumberjack.Logger {
//	return &lumberjack.Logger{
//		Filename:   z.Path,
//		MaxSize:    z.MaxFileSize,
//		MaxBackups: z.MaxBackups,
//		MaxAge:     z.MaxAge,
//		Compress:   z.Compress,
//	}
//}
//
//func (z *zapAdapter) Build() {
//	w := zapcore.AddSync(os.Stdout)
//	if z.FileOut {
//		w = zapcore.NewMultiWriteSyncer(zapcore.AddSync(z.createLumberjackHook()), w)
//	}
//
//	var level zapcore.Level
//	switch z.Level {
//	case "debug":
//		level = zap.DebugLevel
//	case "info":
//		level = zap.InfoLevel
//	case "warn":
//		level = zap.WarnLevel
//	case "error":
//		level = zap.ErrorLevel
//	case "panic":
//		level = zap.PanicLevel
//	default:
//		level = zap.InfoLevel
//	}
//
//	conf := zap.NewProductionEncoderConfig()
//	conf.EncodeTime = zapcore.ISO8601TimeEncoder
//	cnf := zapcore.NewJSONEncoder(conf)
//	core := zapcore.NewCore(cnf, w, level)
//
//	z.logger = zap.New(core)
//	if z.Caller {
//		z.logger = z.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2))
//	}
//
//	if z.Source != "" {
//		z.logger = z.logger.With(zap.Field{
//			Key:    "source",
//			Type:   zapcore.StringType,
//			String: z.Source,
//		})
//	}
//
//	z.sugar = z.logger.Sugar()
//}
//
//func (z *zapAdapter) Debug(args ...interface{}) {
//	z.sugar.Debug(args...)
//}
//
//func (z *zapAdapter) Info(args ...interface{}) {
//	z.sugar.Info(args...)
//}
//
//func (z *zapAdapter) Warn(args ...interface{}) {
//	z.sugar.Warn(args...)
//}
//
//func (z *zapAdapter) Error(args ...interface{}) {
//	z.sugar.Error(args...)
//}
//
//func (z *zapAdapter) DPanic(args ...interface{}) {
//	z.sugar.DPanic(args...)
//}
//
//func (z *zapAdapter) Panic(args ...interface{}) {
//	z.sugar.Panic(args...)
//}
//
//func (z *zapAdapter) Fatal(args ...interface{}) {
//	z.sugar.Fatal(args...)
//}
//
//func (z *zapAdapter) Debugf(template string, args ...interface{}) {
//	z.sugar.Debugf(template, args...)
//}
//
//func (z *zapAdapter) Infof(template string, args ...interface{}) {
//	z.sugar.Infof(template, args...)
//}
//
//func (z *zapAdapter) Warnf(template string, args ...interface{}) {
//	z.sugar.Warnf(template, args...)
//}
//
//func (z *zapAdapter) Errorf(template string, args ...interface{}) {
//	z.sugar.Errorf(template, args...)
//}
//
//func (z *zapAdapter) DPanicf(template string, args ...interface{}) {
//	z.sugar.DPanicf(template, args...)
//}
//
//func (z *zapAdapter) Panicf(template string, args ...interface{}) {
//	z.sugar.Panicf(template, args...)
//}
//
//func (z *zapAdapter) Fatalf(template string, args ...interface{}) {
//	z.sugar.Fatalf(template, args...)
//}
//
//func (z *zapAdapter) Debugw(msg string, keysAndValues ...interface{}) {
//	z.sugar.Debugw(msg, keysAndValues...)
//}
//
//func (z *zapAdapter) Infow(msg string, keysAndValues ...interface{}) {
//	z.sugar.Infow(msg, keysAndValues...)
//}
//
//func (z *zapAdapter) Warnw(msg string, keysAndValues ...interface{}) {
//	z.sugar.Warnw(msg, keysAndValues...)
//}
//
//func (z *zapAdapter) Errorw(msg string, keysAndValues ...interface{}) {
//	z.sugar.Errorw(msg, keysAndValues...)
//}
//
//func (z *zapAdapter) DPanicw(msg string, keysAndValues ...interface{}) {
//	z.sugar.DPanicw(msg, keysAndValues...)
//}
//
//func (z *zapAdapter) Panicw(msg string, keysAndValues ...interface{}) {
//	z.sugar.Panicw(msg, keysAndValues...)
//}
//
//func (z *zapAdapter) Fatalw(msg string, keysAndValues ...interface{}) {
//	z.sugar.Fatalw(msg, keysAndValues...)
//}
