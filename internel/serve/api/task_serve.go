package main

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/handler"
	"dv/internel/serve/api/internal/svc"
	"flag"
	"fmt"
	"os/exec"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/task_serve.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	go openChrome()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	handler.RegisterWSHandlers(server, ctx)
	handler.RegisterH5Handlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.PrintRoutes()
	server.Start()
}

func openChrome() {
	cmd := exec.Command("cmd", "/C", "start", "chrome.exe", "http://127.0.0.1:8888")
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
