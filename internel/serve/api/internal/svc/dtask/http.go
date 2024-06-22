package dtask

import (
	"crypto/tls"
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/util/curl"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

var (
	defaultHeader http.Header
	defaultClient *http.Client
)

func InitDefaultHttp(cfg config.Config) {
	header := http.Header{}
	for k, v := range cfg.HttpConfig.Headers {
		header.Set(k, v)
	}

	defaultClient = &http.Client{Transport: getHttpProxy(cfg.HttpConfig)}
	defaultHeader = header
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

func parseTaskToHttp(t model.Task) (*http.Request, error) {
	var request *http.Request
	var err error
	switch t.Type {
	case model.TypeUrl:
		request, err = http.NewRequest(http.MethodGet, t.Data, nil)
		if err != nil {
			return nil, err
		}
		request.Header = defaultHeader
	case model.TypeCurl:
		_url, header, err := curl.Parse(t.Data)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(http.MethodGet, _url, nil)
		if err != nil {
			return nil, err
		}
		request.Header = header
	case model.TypeProxy:
		request, err = http.NewRequest(http.MethodGet, t.Data, nil)
		if err != nil {
			return nil, err
		}
		var header http.Header
		if err := json.Unmarshal([]byte(t.HeaderJson), &header); err != nil {
			return nil, err
		}
		if len(header) == 0 {
			header = defaultHeader
		}
		request.Header = header
	default:
		return nil, errors.New("type error")

	}

	return request, nil
}
