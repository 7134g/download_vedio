package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	OpenChrome bool
	rest.RestConf

	DB string // 数据库

	HttpConfig

	TaskControlConfig
}

type HttpConfig struct {
	Headers     map[string]string // 默认请求头
	Proxy       string            // 代理地址
	ProxyStatus bool              // 代理是否开启
}

type TaskControlConfig struct {
	WebProxy string // web监听

	TaskConcurrency     int32  // 允许任务并发数
	TaskM3u8Concurrency int32  // m3u8 分片允许的并发数
	SiteLimit           bool   // 是否开启单站点下载限制
	TaskErrorMaxCount   int    // 任务连续最大错误次数
	SaveDir             string // 存储位置
	TaskErrorDuration   uint   // 错误时候休眠多久后重试(秒)
	UseFfmpeg           bool   // 使用ffmpeg进行合并分片
	FfmpegPath          string // ffmpeg程序所在地址

	Concurrency     uint // 并发数
	ConcurrencyM3u8 uint // m3u8 片段并发大小
}
