package zlog

type Option struct {
	Path        string // 文件绝对地址，如：/home/homework/neso/file.log
	Level       string // 日志输出的级别
	MaxFileSize int    // 日志文件大小的最大值，单位(M)
	MaxBackups  int    // 最多保留备份数
	MaxAge      int    // 日志文件保存的时间，单位(天)
	Compress    bool   // 是否压缩
	Caller      bool   // 日志是否需要显示调用位置
	StdOut      bool   // 是否输出到控制台
	FileOut     bool   // 是否需要文件输出
	Source      string // 标志
}

func defaultOption() *Option {
	return &Option{
		Path:   "",
		Level:  "debug",
		Caller: true,
		StdOut: true,
	}
}
