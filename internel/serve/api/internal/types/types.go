// Code generated by goctl. DO NOT EDIT.
package types

type TaskCreateRequest struct {
	Name      string `json:"name"`                          // 任务名字
	VideoType string `json:"video_type,options=[mp4,m3u8]"` // 视频类型
	Type      string `json:"type,options=[url,curl,all]"`   // 任务类型
	Data      string `json:"data"`                          // url 或者 curl
}

type TaskCreateResponse struct {
}

type TaskInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`       // 任务名字
	VideoType string `json:"video_type"` // 视频类型
	Type      string `json:"type"`       // 任务类型
	Data      string `json:"data"`       // url 或者 curl
	Status    uint   `json:"status"`     // 执行状态
	Score     uint   `json:"score"`      // 现在进度
}

type TaskListRequest struct {
	DbQueryList
}

type TaskListResponse struct {
	Total int64      `json:"total"`
	List  []TaskInfo `json:"list"`
}

type TaskUpdateRequest struct {
	ID        uint   `json:"id"`
	Name      string `json:"name,optional"`       // 任务名字
	VideoType string `json:"video_type,optional"` // 视频类型
	Type      string `json:"type,optional"`       // 任务类型
	Data      string `json:"data,optional"`       // url 或者 curl
}

type TaskUpdateResponse struct {
}

type TaskDeleteRequest struct {
	ID uint `json:"id"`
}

type TaskDeleteResponse struct {
}

type TaskRunRequest struct {
	Stop bool `form:"stop,optional"`
}

type TaskRunResponse struct {
	Message string `json:"message"`
}

type TaskStatusRequest struct {
}

type TaskStatusResponse struct {
	Status   bool   `json:"status"`    // 执行状态
	WebProxy string `json:"web_proxy"` // 填写到浏览器代理的地址
}

type GetConfigRequest struct {
}

type GetConfigResponse struct {
	WebProxy          string `json:"web_proxy"`            // web监听
	Concurrency       uint   `json:"concurrency"`          // 并发数
	ConcurrencyM3u8   uint   `json:"concurrency_m_3_u_8"`  // m3u8 片段并发大小
	SaveDir           string `json:"save_dir"`             // 存储位置
	TaskErrorMaxCount uint   `json:"task_error_max_count"` // 任务连续最大错误次数
	TaskErrorDuration uint   `json:"task_error_duration"`  // 错误时候休眠多久后重试(秒)
	UseFfmpeg         bool   `json:"use_ffmpeg"`           // 使用ffmpeg进行合并分片
	FfmpegPath        string `json:"ffmpeg_path"`          // ffmpeg程序所在地址
}

type SetConfigRequest struct {
	WebProxy          string `json:"web_proxy"`            // web监听
	Concurrency       uint   `json:"concurrency"`          // 并发数
	ConcurrencyM3u8   uint   `json:"concurrency_m_3_u_8"`  // m3u8 片段并发大小
	SaveDir           string `json:"save_dir"`             // 存储位置
	TaskErrorMaxCount uint   `json:"task_error_max_count"` // 任务连续最大错误次数
	TaskErrorDuration uint   `json:"task_error_duration"`  // 错误时候休眠多久后重试(秒)
	UseFfmpeg         bool   `json:"use_ffmpeg"`           // 使用ffmpeg进行合并分片
	FfmpegPath        string `json:"ffmpeg_path"`          // ffmpeg程序所在地址
}

type SetConfigResponse struct {
}

type GetCertRequest struct {
	File string `form:"file"`
}

type GetCertResponse struct {
}

type DbQueryList struct {
	Page     int                    `json:"page,range=[0:),default=1,optional"`
	Size     int                    `json:"size,range=(:500],default=10,optional"`
	OrderKey string                 `json:"order_key,optional"`                // 排序字段
	Order    string                 `json:"order,options=[desc,asc],optional"` // 排序逻辑
	Where    map[string]interface{} `json:"where,optional"`
}

type DbQueryListResponse struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
