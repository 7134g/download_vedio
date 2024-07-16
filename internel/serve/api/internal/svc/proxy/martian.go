package proxy

import (
	"bytes"
	"dv/internel/serve/api/internal/dao"
	"dv/internel/serve/api/internal/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/martian"
	"github.com/google/martian/auth"
	"github.com/google/martian/log"
	"github.com/google/martian/mitm"
	"github.com/google/martian/priority"
	"github.com/google/martian/proxyauth"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	monitorAddress = "127.0.0.1:10888" // 监听地址
	taskDB         dao.TaskModel
)

var (
	httpMartian *martian.Proxy // 拦截器全局对象
	certFlag    bool           // 开启自签证书验证
)

var (
	serverProxyUrlParse *url.URL // 解析代理

	serverProxyFlag     bool   // 启用代理
	serverProxy         string // 服务代理地址
	serverProxyUsername string // 用户名
	serverProxyPassword string // 密码
)

func init() {
	log.SetLevel(log.Silent)
}

func OpenCert() {
	certFlag = true
	if err := LoadCert(); err != nil {
		panic(err)
	}
}

func SetServeProxyAddress(address, username, password string) {
	serverProxyFlag = true
	serverProxy = address
	serverProxyUsername = username
	serverProxyPassword = password
}

func SetTaskDb(taskDb dao.TaskModel) {
	taskDB = taskDb
}

func SetMartianAddress(address string) {
	monitorAddress = address
}

func Martian() error {
	httpMartian = martian.NewProxy()
	if certFlag {
		mc, err := mitm.NewConfig(ca, private)
		if err != nil {
			return err
		}
		httpMartian.SetMITM(mc)
	}

	if serverProxyFlag {
		u, err := url.Parse(serverProxy)
		if err != nil {
			return err
		}
		serverProxyUrlParse = u
	}

	group := priority.NewGroup()
	xs := newFirstCheckpoint()
	group.AddRequestModifier(xs, 10)
	group.AddResponseModifier(xs, 10)
	xa := newAuth(proxyauth.NewModifier())
	group.AddRequestModifier(xa, 12)
	group.AddResponseModifier(xa, 12)
	httpMartian.SetRequestModifier(group)
	httpMartian.SetResponseModifier(group)

	fmt.Printf("listen %s, user proxy %s \n", monitorAddress, serverProxy)
	listener, err := net.Listen("tcp", monitorAddress)
	if err != nil {
		return err
	}

	err = httpMartian.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

type firstCheckpoint struct {
}

func newFirstCheckpoint() *firstCheckpoint {
	return &firstCheckpoint{}
}

func (r *firstCheckpoint) ModifyRequest(req *http.Request) error {
	logx.Debugw(
		"url message",
		logx.Field("method", req.Method),
		logx.Field("url", req.URL.String()))
	req.Header.Del("Accept-Encoding")
	HostUrlMap.AddReqUrl(req.URL.String())
	return nil
}

func (r *firstCheckpoint) ModifyResponse(res *http.Response) error {
	reqUrl := res.Request.URL.String()
	HostUrlMap.AddReqUrl(reqUrl)

	parts := strings.Split(res.Request.URL.Path, ".")
	if len(parts) > 0 {
		ext := parts[len(parts)-1]
		switch ext {
		case model.VideoTypeMp4, model.VideoTypeM3u8:
			v, _ := json.Marshal(res.Request.Header)
			header := string(v)
			findTask, _ := taskDB.Exist(reqUrl)
			if findTask == nil {
				t := &model.Task{
					Name:       fmt.Sprintf("%s", time.Now().Format("2006_01_02_15_04_05")),
					VideoType:  ext,
					Type:       model.TypeProxy,
					Url:        reqUrl,
					HeaderJson: header,
				}
				if err := taskDB.Create(t); err != nil {
					logx.Error(err)
				}

				t, _ = taskDB.Exist(reqUrl)
				sourceChan <- message{
					taskId: t.ID,
					source: t.Url,
				}
			}
		case "html":
			data, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			if len(data) == 0 {
				return nil
			}

			HostUrlMap.AddBody(reqUrl, data)
			res.Body = io.NopCloser(bytes.NewBuffer(data))
		default:
			return nil
		}

	}

	return nil
}

type xauth struct {
	pAuth *proxyauth.Modifier
}

func newAuth(pAuth *proxyauth.Modifier) *xauth {
	return &xauth{pAuth: pAuth}
}

func (r *xauth) ModifyRequest(req *http.Request) error {
	if serverProxy == "" {
		return nil
	}

	httpMartian.SetDownstreamProxy(serverProxyUrlParse)

	if serverProxyUsername != "" {
		un := base64.StdEncoding.EncodeToString([]byte(serverProxyUsername))
		pw := base64.StdEncoding.EncodeToString([]byte(serverProxyPassword))
		//req.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s:%s", un, pw))
		ctx := martian.NewContext(req)
		authCTX := auth.FromContext(ctx)
		if authCTX.ID() != fmt.Sprintf("%s:%s", un, pw) {
			authCTX.SetError(errors.New("auth error"))
			ctx.SkipRoundTrip()
		}
	}

	return nil
}

func (r *xauth) ModifyResponse(res *http.Response) error {
	if serverProxy == "" {
		return nil
	}
	return r.pAuth.ModifyResponse(res)
}

func getNumber() int64 {
	currentTime := time.Now().Unix()
	return currentTime - currentTime%10
}
