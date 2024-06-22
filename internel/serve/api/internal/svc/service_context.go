package svc

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/dao"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/middleware"
	"dv/internel/serve/api/internal/svc/dtask"
	"dv/internel/serve/api/internal/svc/proxy"
	"dv/internel/serve/api/internal/util/files"
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
	TaskModel       dao.TaskModel
	LogData         *logCache

	Hub *ws_conn.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	// db
	dao.DaoInit(db.InitSqlite(c.DB))
	taskModel := dao.TaskDao
	dtask.Control.Init(c)

	// 开启被动代理
	threading.GoSafe(func() {
		proxy.SetTaskDb(taskModel)
		proxy.SetServeProxyAddress(c.Proxy, "", "")
		proxy.OpenCert()
		proxy.SetMartianAddress(c.WebProxy)
		if err := proxy.Martian(); err != nil {
			panic(err)
		}
	})
	// 处理 ProxyCatchUrl 和 ProxyCatchHtml 匹配
	threading.GoSafe(func() {
		proxy.MatchInformation()
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
		LogData:         logData,
		Hub:             hub,
	}
}
