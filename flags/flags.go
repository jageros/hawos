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

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

var (
	Options *Option // 全局存储
	v       *viper.Viper
)

type ValInfo struct {
	Val         interface{}
	Description string
}

// 命令行启动参数和配置文件统一的key，避免多次书写出错，命令行中直接使用即可，配置文件中，点隔开代表分级
var (
	//config
	keyConfig = "config"
	//server
	keyId   = "server.id"
	keyName = "server.name"
	keyMode = "server.mode"
	// log
	keyLogDir         = "log.dir"
	keyLogCaller      = "log.caller"
	keyLogStdout      = "log.stdout"
	keyLogMaxFileSize = "log.max_file_size"
	keyLogMaxBackups  = "log.max_backups"
	keyLogMaxAge      = "log.max_age"
	keyLogCompress    = "log.compress"
)

// Option 配置数据结构体
type Option struct {
	ID         int    // 服务id
	AppName    string // 服务名称
	Mode       string // 模式
	Configfile string // 配置文件路径

	// log配置
	LogDir          string // log目录
	LogCaller       bool   // 是否开启记录输出日志的代码文件行号
	LogStdout       bool   // 是否输出到控制台
	LogMaxFileSize  int    // 最大日志文件大小 MB
	LogMaxBackups   int    // 最大备份数量
	LogMaxAge       int    // 最大日志天数
	LogFileCompress bool   // 是否压缩日志

	// other conf
	Keys map[string]*ValInfo

	OnReload func()
}

// AppID 返回字符串型的appid
func AppID() string {
	return strconv.Itoa(Options.ID)
}

// Source 可用于服务注册中的key或者日志中的source可区别开不同的服务或同个服务不同的结点
func Source() string {
	return fmt.Sprintf("%s/%d", Options.AppName, Options.ID)
}

// defaultOption 返回程序默认的配置项数据
func defaultOption(name string) *Option {
	op := &Option{
		ID:        1,
		AppName:   name,
		Mode:      "debug",
		LogStdout: true,
	}
	return op
}

// load 从viper中获取解析出来的参数初始化option中的字段
func (op *Option) load(v *viper.Viper) {
	//server
	op.ID = v.GetInt(keyId)
	op.AppName = v.GetString(keyName)
	op.Mode = v.GetString(keyMode)
	//log
	op.LogDir = v.GetString(keyLogDir)
	op.LogCaller = v.GetBool(keyLogCaller)
	op.LogStdout = v.GetBool(keyLogStdout)
	op.LogMaxFileSize = v.GetInt(keyLogMaxFileSize)
	op.LogMaxBackups = v.GetInt(keyLogMaxBackups)
	op.LogMaxAge = v.GetInt(keyLogMaxAge)
	op.LogFileCompress = v.GetBool(keyLogCompress)

	err := logx.Init(func(opt *logx.Option) {
		opt.Level = op.Mode
		opt.LogPath = op.LogDir
		opt.Stdout = op.LogStdout
		opt.Source = op.AppName + strconv.Itoa(op.ID)
		opt.Caller = op.LogCaller
		opt.Stdout = op.LogStdout
		opt.MaxFileSize = op.LogMaxFileSize
		opt.MaxBackups = op.LogMaxBackups
		opt.MaxAge = op.LogMaxAge
		opt.Compress = op.LogFileCompress
	})

	if err != nil {
		log.Fatalf("logx init err: %v", err)
	}
}

// Parse 解析配置， 启动参数有传参则忽略配置文件
func Parse(name string, opts ...func(opt *Option)) (ctx contextx.Context, wait func(), cancel func()) {
	Options = defaultOption(name)

	// 调用该接口时可以改变默认值，但优先顺序为 启动参数 > 配置文件 > 接口传参 > 默认值
	for _, optf := range opts {
		optf(Options)
	}

	// 启动参数：
	//config
	pflag.String(keyConfig, Options.Configfile, "Config file path")
	// server
	pflag.Int(keyId, Options.ID, "Application Id")
	pflag.String(keyName, Options.AppName, "Application name")
	pflag.String(keyMode, Options.Mode, "Server mode, default: debug, optional：debug/test/release")
	// log
	pflag.String(keyLogDir, Options.LogDir, "Log dir")
	pflag.Bool(keyLogCaller, Options.LogCaller, "log caller")
	pflag.Bool(keyLogStdout, Options.LogStdout, "log stdout")

	// other
	for key, v := range Options.Keys {
		val := v.Val
		usage := v.Description
		if usage == "" {
			usage = key
		}
		switch val.(type) {
		case string:
			pflag.String(key, val.(string), usage)
		case []string:
			pflag.StringSlice(key, val.([]string), usage)
		case uint, int, uint8, int8, uint16, int16, uint32, int32, uint64, int64:
			pflag.Int(key, utils.ToInt(val), usage)
		case []uint, []int, []uint8, []int8, []uint16, []int16, []uint32, []int32, []uint64, []int64:
			pflag.IntSlice(key, utils.ToIntSlice(val), usage)
		case float32, float64:
			pflag.Float64(key, utils.ToFloat64(val), usage)
		case []float32, []float64:
			pflag.Float64Slice(key, utils.ToFloat64Slice(val), usage)
		case bool:
			pflag.Bool(key, val.(bool), usage)
		case []bool:
			pflag.BoolSlice(key, val.([]bool), usage)
		}
	}

	pflag.Parse()

	v = viper.New()
	err := v.BindPFlags(pflag.CommandLine) // 绑定命令行参数
	if err != nil {
		log.Fatalf("v.BindPFlags err: %v", err)
	}

	// 获取命令行配置文件路径参数
	path := v.GetString(keyConfig)

	var dir, fileName, fileType string
	// 配置文件路径不为空则读取配置文件， 即命令行有传入配置文件路径时
	if path != "" {
		// 获取配置文件的后缀名
		strs := strings.Split(path, ".")
		if len(strs) < 2 {
			log.Fatal("错误的配置文件路径")
		}

		fileType = strs[len(strs)-1]
		switch fileType {
		// 支持的文件类型yaml、json、 toml、hcl, ini
		case "yaml", "yml", "json", "toml", "hcl", "ini":
			fPath := strings.Replace(path, "."+fileType, "", -1)
			strs = strings.Split(fPath, "/")
			if len(strs) < 2 {
				fileName = strs[0]
				dir = "./"
			} else {
				fileName = strs[len(strs)-1]
				dir = strings.Join(strs[:len(strs)-1], "/")
			}

			//设置读取的配置文件
			v.SetConfigName(fileName)
			//添加读取的配置文件路径
			v.AddConfigPath(dir)
			//设置配置文件类型
			v.SetConfigType(fileType)

			if err := v.ReadInConfig(); err != nil {
				log.Fatalf("v.ReadInConfig err: %v", err)
			}

			// 这部分代码为监听配置文件是否有更新，有更新则重新解析配置文件，重新解析也无法覆盖命令行参数
			// 可根据需求启用
			// 依赖库： "github.com/fsnotify/fsnotify"
			//设置监听回调函数
			v.OnConfigChange(func(e fsnotify.Event) {
				if e.Op == fsnotify.Write {
					logx.Sync()
					Options.load(v)
					if Options.OnReload != nil {
						Options.OnReload()
					}
				}
			})
			//开始监听
			v.WatchConfig()

		default:
			// 其他类型抛出错误
			log.Fatal("错误的配置文件类型")
		}
	}

	// 从viper解析出来的参数对option中的字段赋值
	Options.load(v)

	// 结合信号和context实现的类似与errgroup的一个库，可以根据自己项目需求设计自己的接口
	ctx, cancel = contextx.Default()

	wait = func() {
		err := ctx.Wait()
		logx.Err(err).Msg("Application Stop!")
		logx.Sync()
	}

	return
}

func GetString(key string) string {
	return v.GetString(key)
}

func GetStringSlice(key string) []string {
	return v.GetStringSlice(key)
}

func GetInt(key string) int {
	return v.GetInt(key)
}

func GetIntSlice(key string) []int {
	return v.GetIntSlice(key)
}

func GetFloat64(key string) float64 {
	return v.GetFloat64(key)
}

func GetBool(key string) bool {
	return v.GetBool(key)
}

func GetVal(key string) interface{} {
	return v.Get(key)
}
