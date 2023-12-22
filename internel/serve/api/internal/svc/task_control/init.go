package task_control

import (
	"context"
	"crypto/tls"
	"dv/internel/serve/api/internal/config"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var (
	tc       *TaskControl
	tcConfig taskControlConfig
)

type TaskControl struct {
	wg      sync.WaitGroup
	mux     sync.Mutex
	running bool
	ctx     context.Context
	cancel  context.CancelFunc
	vacancy chan struct{} // 并发控制
}

type taskControlConfig struct {
	// http
	Transport http.RoundTripper
	Headers   http.Header

	// 执行设置
	cfg config.TaskControlConfig
}

func NewTaskControl(c config.Config) *TaskControl {
	ctx, cancel := context.WithCancel(context.Background())
	tc = &TaskControl{
		wg:      sync.WaitGroup{},
		mux:     sync.Mutex{},
		running: false,
		ctx:     ctx,
		cancel:  cancel,
		vacancy: make(chan struct{}, c.TaskControlConfig.Concurrency),
	}
	tcConfig = taskControlConfig{
		Transport: getHttpProxy(c.HttpConfig),
		Headers:   getHeader(c.HttpConfig),
		cfg:       c.TaskControlConfig,
	}

	return tc
}

func getHttpProxy(c config.HttpConfig) http.RoundTripper {
	if !c.ProxyStatus {
		return nil
	}

	httpProxy := func(proxy string) func(*http.Request) (*url.URL, error) {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Fatalln(err)
		}
		return http.ProxyURL(proxyUrl)
	}

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, // 使用环境变量的代理
		Proxy: httpProxy(c.Proxy),
	}

	return t
}

func getHeader(c config.HttpConfig) http.Header {
	header := http.Header{}
	for k, v := range c.Headers {
		header.Set(k, v)
	}

	return header
}
