package svc

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/middleware"
	proxy2 "dv/internel/serve/api/internal/svc/proxy"
	"dv/internel/serve/api/internal/svc/task_control"
	"dv/internel/serve/api/internal/util/files"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/ws_conn"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest"
	"io"
	"os"
	"time"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware
	TaskModel       *model.TaskModel
	TaskControl     *task_control.TaskControl
	LogData         *logCache

	Hub *ws_conn.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	// db
	taskModel := model.NewTaskModel(db.InitSqlite(c.DB))
	task_control.InitTask(c)

	// 开启被动代理
	threading.GoSafe(func() {
		proxy2.SetTaskDb(taskModel)
		proxy2.SetServeProxyAddress(c.Proxy, "", "")
		proxy2.OpenCert()
		proxy2.SetMartianAddress(c.WebProxy)
		if err := proxy2.Martian(); err != nil {
			panic(err)
		}
	})
	// 处理 ProxyCatchUrl 和 ProxyCatchHtml 匹配
	threading.GoSafe(func() {
		proxy2.MatchInformation()
	})

	// 设置日志
	f, err := files.GetFile(fmt.Sprintf("./log/%s.log", time.Now().Format(time.DateOnly)))
	if err != nil {
		panic(err)
	}
	logData := newLogCache()
	logWrite := logx.NewWriter(io.MultiWriter(os.Stdout, logData, f))
	logx.SetWriter(logWrite)

	// ws
	hub := ws_conn.NewHub()
	threading.GoSafe(hub.Run)

	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
		TaskModel:       taskModel,
		TaskControl:     task_control.NewTaskControl(c.TaskControlConfig.Concurrency),
		LogData:         logData,
		Hub:             hub,
	}
}
