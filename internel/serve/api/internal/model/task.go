package model

const (
	VideoTypeMp4        = "mp4"
	VideoTypeM3u8       = "m3u8"
	VideoTypeM3u8Single = "m3u8_single"
)

const (
	TypeUrl   = "url"
	TypeCurl  = "curl"
	TypeProxy = "proxy"
)

const (
	StatusWait = iota
	StatusRunning
	StatusError
	StatusSuccess
	StatusStop
)

type Task struct {
	ID         int    `json:"id" gorm:"primaryKey;column:id"`
	Name       string `json:"name"`        // 任务名字
	VideoType  string `json:"video_type"`  // 视频类型
	Type       string `json:"type"`        // 任务类型
	Url        string `json:"data"`        // url
	HeaderJson string `json:"header_json"` // 序列化的请求头
	Status     int    `json:"status"`      // 执行状态
}
