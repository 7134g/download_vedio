// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	task "dv/internel/serve/api/internal/handler/task"
	"dv/internel/serve/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthInterceptor},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/create",
					Handler: task.CreateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/list",
					Handler: task.ListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/update",
					Handler: task.UpdateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/delete",
					Handler: task.DeleteHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/run",
					Handler: task.RunHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/status",
					Handler: task.StatusHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/config",
					Handler: task.GetConfigHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/config/set",
					Handler: task.SetConfigHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/task"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/get_file",
				Handler: GetCertFileHandler(serverCtx),
			},
		},
	)
}
