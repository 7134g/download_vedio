package proxy

import (
	"bytes"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// 用于读出 response 再重新写入
type writer struct {
	write        http.ResponseWriter
	code         int
	responseBody *bytes.Buffer
}

func newWrite(write http.ResponseWriter) *writer {
	return &writer{
		write:        write,
		responseBody: bytes.NewBuffer(nil),
	}
}

func (w *writer) Header() http.Header {
	return w.write.Header()
}

func (w *writer) Write(bytes []byte) (int, error) {
	w.responseBody.Write(bytes)
	return w.write.Write(bytes)
}

func (w *writer) WriteHeader(statusCode int) {
	w.code = statusCode
	w.write.WriteHeader(statusCode)
}

func ParseHtmlTitle(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	node := doc.Find("title")
	if node == nil {
		return "", errors.New("cannot find title")
	}

	fn := node.First()
	if fn == nil {
		return "", errors.New("cannot find title")
	}
	title := fn.Text()
	return title, nil
}

// ExtractRequestToString 提取请求包
func ExtractRequestToString(res *http.Request) string {
	buf := bytes.NewBuffer([]byte{})
	defer buf.Reset()
	err := res.Write(buf)
	if err != nil {
		return ""
	}

	return buf.String()
}

var (
	regUrl, _  = regexp.Compile(`([^\/]+)(\.m3u8|\.mp4)$`)
	tickerTime = time.Second * 10
)

func MatchInformation() {
	ticker := time.NewTicker(tickerTime)

	for {

		deleteKey := []string{}

		select {
		case <-ticker.C:
			table.ProxyCatchUrl.Each(func(link string, taskId uint) {
				// https://1257120875.vod2.myqcloud.com/0ef121cdvodtransgzp1257120875/3055695e5285890780828799271/v.f230.m3u8
				u, err := url.Parse(link)
				if err != nil {
					return
				}
				var filename, name string
				filename = regUrl.FindString(u.Path)
				parts := strings.Split(filename, ".")
				if len(parts) > 1 {
					name = parts[0]
				}

				table.ProxyCatchHtmlTitle.Each(func(title string, html string) {
					regKey, _ := regexp.Compile(fmt.Sprintf("(%s)|(%s)|(%s)", link, filename, name))
					if regKey.MatchString(html) {
						logx.Debugf("taskId %d change name %s", taskId, title)
						logx.Info("添加任务：", title)
						if err := taskDB.Update(&model.Task{ID: taskId, Name: title}); err != nil {
							logx.Error(err)
						} else {
							deleteKey = append(deleteKey, title, link)
						}
					}

				})
			})

		}

		for _, k := range deleteKey {
			table.ProxyCatchUrl.Del(k)
			table.ProxyCatchHtmlTitle.Del(k)
		}

	}
}
