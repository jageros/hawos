/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    flag
 * @Date:    2021/6/18 4:36 下午
 * @package: flag
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package flags
//
//import (
//	"fmt"
//	"github.com/jageros/hawox/attribute"
//	"github.com/jageros/hawox/contextx"
//	"github.com/jageros/hawox/logx"
//	"github.com/jageros/hawox/mysql"
//	"github.com/jageros/hawox/redis"
//	"github.com/spf13/pflag"
//	"github.com/spf13/viper"
//	logx2 "github.com/tal-tech/go-zero/core/logx"
//	clientv3 "go.etcd.io/etcd/client/v3"
//	"log"
//	"strconv"
//	"strings"
//)
//
//var (
//	Options *Option // 全局存储
//	v       *viper.Viper
//)
//
//// 命令行启动参数和配置文件统一的key，避免多次书写出错，命令行中直接使用即可，配置文件中，点隔开代表分级
//var (
//	//config
//	keyConfig = "config"
//	//server
//	keyId   = "server.id"
//	keyMode = "server.mode"
//	// log
//	keyLogDir     = "log.dir"
//	keyLogCaller  = "log.caller"
//	keyLogRequest = "log.request"
//	keyLogStat    = "log.stat"
//	// http
//	keyHttpIp   = "http.ip"
//	keyHttpPort = "http.port"
//	// rpc
//	keyRpcIp   = "rpc.ip"
//	keyRpcPort = "rpc.port"
//	// mysql
//	keyMysqlAddr = "mysql.addr"
//	keyMysqlUser = "mysql.user"
//	keyMysqlPwd  = "mysql.password"
//	keyMysqlDB   = "mysql.database"
//	// mongo
//	keyMongoAddr = "mongo.addr"
//	keyMongoUser = "mongo.user"
//	keyMongoPwd  = "mongo.password"
//	keyMongoDB   = "mongo.database"
//	// etcd
//	keyEtcdAddr     = "etcd.addrs"
//	keyEtcdUser     = "etcd.user"
//	keyEtcdPassword = "etcd.password"
//	//redis
//	keyRedisAddr     = "redis.addrs"
//	keyRedisDB       = "redis.db"
//	keyRedisUser     = "redis.user"
//	keyRedisPassword = "redis.password"
//	//queue
//	keyQueueType = "queue.type"
//	keyQueueAddr = "queue.addrs"
//	//frontend
//	keyFrontendAddr = "frontend.addr"
//)
//
//// Option 配置数据结构体
//type Option struct {
//	ID         int    // 服务id
//	AppName    string // 服务名称
//	Mode       string // 模式
//	Configfile string // 配置文件路径
//
//	// log配置
//	LogDir     string // log目录
//	LogCaller  bool   // 是否开启记录输出日志的代码文件行号
//	LogRequest bool   // 是否开启http请求的日志记录
//	LogStat    bool   // 是否开启系统状态日志输出
//
//	// http服务配置
//	HttpIp   string
//	HttpPort int
//
//	// rpc服务配置
//	RpcIp   string
//	RpcPort int
//
//	// mysql 配置
//	MysqlAddr string
//	MysqlUser string
//	MysqlPwd  string
//	MysqlDB   string
//
//	// MongoDB 配置
//	MongoAddr string
//	MongoUser string
//	MongoPwd  string
//	MongoDB   string
//
//	// etcd 配置
//	EtcdAddrs  string // 地址集，用;隔开
//	EtcdUser   string
//	EtcdPasswd string
//
//	// redis配置
//	RedisAddrs  string // 数据库地址集，用;隔开
//	RedisDBNo   int    // 数据库编号
//	RedisUser   string
//	RedisPasswd string
//
//	// 消息队列配置
//	QueueType  string // 消息队列类型
//	QueueAddrs string // 地址
//
//	// websocket注册进etcd的地址（外网地址），将返回给客户端
//	FrontendAddr string
//}
//
//// AppID 返回字符串型的appid
//func AppID() string {
//	return strconv.Itoa(Options.ID)
//}
//
//// Source 可用于服务注册中的key或者日志中的source可区别开不同的服务或同个服务不同的结点
//func Source() string {
//	return fmt.Sprintf("%s/%d", Options.AppName, Options.ID)
//}
//
//// defaultOption 返回程序默认的配置项数据
//func defaultOption() *Option {
//	op := &Option{
//		ID:       1,
//		Mode:     "debug",
//		HttpIp:   "127.0.0.1",
//		HttpPort: 8888,
//		RpcIp:    "127.0.0.1",
//		RpcPort:  9999,
//	}
//	return op
//}
//
////// logPath log文件输出的路径，根据AppName来命名
////func logPath() string {
////	if strings.HasSuffix(Options.LogDir, "/") {
////		return fmt.Sprintf("%s%s.log", Options.LogDir, Options.AppName)
////	}
////	return fmt.Sprintf("%s/%s.log", Options.LogDir, Options.AppName)
////}
//
//// load 从viper中获取解析出来的参数初始化option中的字段
//func (op *Option) load(v *viper.Viper) {
//	//server
//	op.ID = v.GetInt(keyId)
//	op.Mode = v.GetString(keyMode)
//	//log
//	op.LogDir = v.GetString(keyLogDir)
//	op.LogCaller = v.GetBool(keyLogCaller)
//	op.LogRequest = v.GetBool(keyLogRequest)
//	op.LogStat = v.GetBool(keyLogStat)
//	//http
//	op.HttpIp = v.GetString(keyHttpIp)
//	op.HttpPort = v.GetInt(keyHttpPort)
//	//rpc
//	op.RpcIp = v.GetString(keyRpcIp)
//	op.RpcPort = v.GetInt(keyRpcPort)
//	//mysql
//	op.MysqlAddr = v.GetString(keyMysqlAddr)
//	op.MysqlUser = v.GetString(keyMysqlUser)
//	op.MysqlPwd = v.GetString(keyMysqlPwd)
//	op.MysqlDB = v.GetString(keyMysqlDB)
//	//mongo
//	op.MongoAddr = v.GetString(keyMongoAddr)
//	op.MongoUser = v.GetString(keyMongoUser)
//	op.MongoPwd = v.GetString(keyMongoPwd)
//	op.MongoDB = v.GetString(keyMongoDB)
//	//etcd
//	op.EtcdAddrs = v.GetString(keyEtcdAddr)
//	op.EtcdUser = v.GetString(keyEtcdUser)
//	op.EtcdPasswd = v.GetString(keyEtcdPassword)
//	//redis
//	op.RedisAddrs = v.GetString(keyRedisAddr)
//	op.RedisUser = v.GetString(keyRedisUser)
//	op.RedisPasswd = v.GetString(keyRedisPassword)
//	//queue
//	op.QueueType = v.GetString(keyQueueType)
//	op.QueueAddrs = v.GetString(keyQueueAddr)
//	//frontend
//	op.FrontendAddr = v.GetString(keyFrontendAddr)
//	if op.FrontendAddr == "" {
//		op.FrontendAddr = fmt.Sprintf("127.0.0.1:%d", op.HttpPort)
//	}
//
//	// 根据配置中的mode设置log的等级
//	var logLevel string
//	switch op.Mode {
//	case "release":
//		logLevel = logx.InfoLevel
//	case "info":
//		logLevel = logx.InfoLevel
//	case "warn":
//		logLevel = logx.WarnLevel
//	case "error":
//		logLevel = logx.ErrorLevel
//	case "panic":
//		logLevel = logx.PanicLevel
//	default:
//		logLevel = logx.DebugLevel
//	}
//
//	// 配置日志
//	var logOpfs []logx.Option
//	if op.LogCaller {
//		logOpfs = append(logOpfs, logx.SetCaller())
//	}
//	if op.LogRequest {
//		logOpfs = append(logOpfs, logx.SetRequest())
//	}
//	source := Source()
//	if source != "" {
//		logOpfs = append(logOpfs, logx.SetSource(source))
//	}
//	if op.LogDir != "" {
//		logOpfs = append(logOpfs, logx.SetFileOut(op.LogDir, op.AppName))
//	}
//
//	if op.LogStat {
//		logx.SetupZeroLog(logx2.LogConf{
//			ServiceName: op.AppName,
//			Mode:        "file",
//			Level:       "info",
//			Path:        op.LogDir + "/stat",
//			KeepDays:    7,
//		})
//	} else {
//		logx.DisableZeroLog()
//	}
//
//	// 初始化日志
//	logx.Init(logLevel, logOpfs...)
//
//}
//
//func DefaultRedis(option *Option) {
//	option.RedisAddrs = "127.0.0.1:6379"
//}
//
//func DefaultClusterRedis(addrs ...string) func(*Option) {
//	addr := "127.0.0.1:7001;127.0.0.1:7002;127.0.0.1:7003;127.0.0.1:7004;127.0.0.1:7005;127.0.0.1:7006"
//	if len(addrs) > 0 {
//		addr = strings.Join(addrs, ";")
//	}
//	return func(option *Option) {
//		option.RedisAddrs = addr
//	}
//}
//
//func DefaultMysql(option *Option) {
//	option.MysqlAddr = "127.0.0.1:3306"
//	option.MysqlUser = "root"
//	option.MysqlPwd = "QianYin@66"
//	option.MysqlDB = "hawox"
//}
//
//func DefaultMongo(option *Option) {
//	option.MongoAddr = "127.0.0.1:27017"
//	option.MongoDB = "attribute"
//}
//
//func DefaultEtcd(option *Option) {
//	option.EtcdAddrs = "127.0.0.1:2379"
//}
//
//func DefaultQueue(option *Option) {
//	option.QueueType = "kafka"
//	option.QueueAddrs = "127.0.0.1:9092"
//}
//
//// Parse 解析配置， 启动参数有传参则忽略配置文件
//func Parse(name string, opts ...func(*Option)) (ctx contextx.Context, wait func()) {
//	Options = defaultOption()
//
//	// 调用该接口时可以改变默认值，但优先顺序为 启动参数 > 配置文件 > 接口传参 > 默认值
//	for _, optf := range opts {
//		optf(Options)
//	}
//
//	// 启动参数：
//	//config
//	pflag.String(keyConfig, Options.Configfile, "Config file path")
//	// server
//	pflag.Int(keyId, Options.ID, "Application Id")
//	pflag.String(keyMode, Options.Mode, "Server mode, default: debug, optional：debug/test/release")
//	// log
//	pflag.String(keyLogDir, Options.LogDir, "Log dir")
//	pflag.Bool(keyLogCaller, Options.LogCaller, "log caller")
//	pflag.Bool(keyLogRequest, Options.LogRequest, "log request")
//	pflag.Bool(keyLogStat, Options.LogStat, "log stat")
//	// http
//	pflag.String(keyHttpIp, Options.HttpIp, "Http listen ip")
//	pflag.Int(keyHttpPort, Options.HttpPort, "HTTP server port")
//	// rpc
//	pflag.String(keyRpcIp, Options.RpcIp, "Rpc listen ip")
//	pflag.Int(keyRpcPort, Options.RpcPort, "Rpc server port")
//	// mysql
//	pflag.String(keyMysqlAddr, Options.MysqlAddr, "mysql addr")
//	pflag.String(keyMysqlUser, Options.MysqlUser, "mysql user")
//	pflag.String(keyMysqlPwd, Options.MysqlPwd, "mysql password")
//	pflag.String(keyMysqlDB, Options.MysqlDB, "mysql database")
//	// mongo
//	pflag.String(keyMongoAddr, Options.MongoAddr, "MongoDB addr")
//	pflag.String(keyMongoUser, Options.MongoUser, "MongoDB user")
//	pflag.String(keyMongoPwd, Options.MongoPwd, "MongoDB password")
//	pflag.String(keyMongoDB, Options.MongoDB, "MongoDB database")
//	// etcd
//	pflag.String(keyEtcdAddr, Options.EtcdAddrs, "Etcd addrs, use ; split")
//	pflag.String(keyEtcdUser, Options.EtcdUser, "Etcd user")
//	pflag.String(keyEtcdPassword, Options.EtcdPasswd, "Etcd password")
//	// redis
//	pflag.String(keyRedisAddr, Options.RedisAddrs, "Redis addrs, use ; split")
//	pflag.Int(keyRedisDB, Options.RedisDBNo, "Redis DB NO.")
//	pflag.String(keyRedisUser, Options.RedisUser, "Redis user")
//	pflag.String(keyRedisPassword, Options.RedisPasswd, "Redis password")
//	// queue
//	pflag.String(keyQueueType, Options.QueueType, "queue type, optional: nsq/kafka")
//	pflag.String(keyQueueAddr, Options.QueueAddrs, "queue addrs, use ; split")
//	// frontend
//	pflag.String(keyFrontendAddr, Options.FrontendAddr, "frontend addr")
//
//	pflag.Parse()
//
//	v = viper.New()
//	err := v.BindPFlags(pflag.CommandLine) // 绑定命令行参数
//	if err != nil {
//		log.Fatalf("v.BindPFlags err: %v", err)
//	}
//
//	// 获取命令行配置文件路径参数
//	path := v.GetString(keyConfig)
//
//	var dir, fileName, fileType string
//	// 配置文件路径不为空则读取配置文件， 即命令行有传入配置文件路径时
//	if path != "" {
//		// 获取配置文件的后缀名
//		strs := strings.Split(path, ".")
//		if len(strs) < 2 {
//			log.Fatal("错误的配置文件路径")
//		}
//
//		fileType = strs[len(strs)-1]
//		switch fileType {
//		// 支持的文件类型yaml、json、 toml、hcl, ini
//		case "yaml", "yml", "json", "toml", "hcl", "ini":
//			fPath := strings.Replace(path, "."+fileType, "", -1)
//			strs = strings.Split(fPath, "/")
//			if len(strs) < 2 {
//				fileName = strs[0]
//				dir = "./"
//			} else {
//				fileName = strs[len(strs)-1]
//				dir = strings.Join(strs[:len(strs)-1], "/")
//			}
//
//			//设置读取的配置文件
//			v.SetConfigName(fileName)
//			//添加读取的配置文件路径
//			v.AddConfigPath(dir)
//			//设置配置文件类型
//			v.SetConfigType(fileType)
//
//			if err := v.ReadInConfig(); err != nil {
//				log.Fatalf("v.ReadInConfig err: %v", err)
//			}
//
//			// 这部分代码为监听配置文件是否有更新，有更新则重新解析配置文件，重新解析也无法覆盖命令行参数
//			// 可根据需求启用
//			// 依赖库： "github.com/fsnotify/fsnotify"
//			////设置监听回调函数
//			//v.OnConfigChange(func(e fsnotify.Event) {
//			//  logx.Sync() // 重新初始化时先同步
//			//	Options.load(v)
//			//	if onReloadConfigFn != nil {
//			//		onReloadConfigFn()
//			//	}
//			//})
//			////开始监听
//			//v.WatchConfig()
//
//		default:
//			// 其他类型抛出错误
//			log.Fatal("错误的配置文件类型")
//		}
//	}
//
//	Options.AppName = name
//	// 从viper解析出来的参数对option中的字段赋值
//	Options.load(v)
//
//	// 结合信号和context实现的类似与errgroup的一个库，可以根据自己项目需求设计自己的接口
//	// 不可直接照搬，这里存在其他依赖库暂时未开源，所以属于伪代码
//	ctx, _ = contextx.Default()
//
//	wait = func() {
//		err := ctx.Wait()
//		logx.Infof("Application Stop With: %v", err)
//		logx.Sync()
//	}
//
//	return
//}
//
//// ===================
//
//func RedisConfigOpf(rdb redis.RdbConfig) {
//	if len(Options.RedisAddrs) <= 0 {
//		DefaultRedis(Options)
//	}
//	rdb.SetAddrs(Options.RedisAddrs)
//	rdb.SetUsername(Options.RedisUser)
//	rdb.SetPassword(Options.RedisPasswd)
//	rdb.SetDB(Options.RedisDBNo)
//}
//
//func EtcdConfigOpf(cfg *clientv3.Config) {
//	if len(Options.EtcdAddrs) <= 0 {
//		DefaultEtcd(Options)
//	}
//	addrs := strings.Split(Options.EtcdAddrs, ";")
//	cfg.Endpoints = addrs
//	cfg.Username = Options.EtcdUser
//	cfg.Password = Options.EtcdPasswd
//}
//
//func MysqlConfigOpf(opt *mysql.Option) {
//	if len(Options.MysqlAddr) <= 0 {
//		DefaultMysql(Options)
//	}
//	opt.Addr = Options.MysqlAddr
//	opt.Username = Options.MysqlUser
//	opt.Password = Options.MysqlPwd
//	opt.Database = Options.MysqlDB
//}
//
//func AttributeConfigOpf(opt *attribute.Option) {
//	if len(Options.MongoAddr) <= 0 {
//		DefaultMongo(Options)
//	}
//	opt.Addr = Options.MongoAddr
//	opt.DBName = Options.MongoDB
//	opt.User = Options.MongoUser
//	opt.Password = Options.MongoPwd
//}
//
//func GetInt(key string) int {
//	return v.GetInt(key)
//}
//
//func GetIntSlice(key string) []int {
//	return v.GetIntSlice(key)
//}
//
//func GetFloat64(key string) float64 {
//	return v.GetFloat64(key)
//}
//
//func GetBool(key string) bool {
//	return v.GetBool(key)
//}
//
//func GetString(key string) string {
//	return v.GetString(key)
//}
//
//func GetStringSlice(key string) []string {
//	return v.GetStringSlice(key)
//}
//
//func GetVal(key string) interface{} {
//	return v.Get(key)
//}
