package handler

import (
	"dv/internel/serve/api/internal/handler/h5"
	"dv/internel/serve/api/internal/svc"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterH5Handlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/index.css",
					Handler: h5.Css,
				},
				{
					Method:  http.MethodGet,
					Path:    "/index.js",
					Handler: h5.Js,
				},
				{
					Method:  http.MethodGet,
					Path:    "/",
					Handler: h5.Html,
				},
			}...,
		),
		rest.WithPrefix("/"),
	)
}
