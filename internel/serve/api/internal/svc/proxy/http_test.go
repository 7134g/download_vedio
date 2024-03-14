package proxy

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestNewWrite(t *testing.T) {
	response := &httptest.ResponseRecorder{}
	lw := newWrite(response)
	_, _ = lw.Write([]byte("abc"))
	t.Log(lw.responseBody.String())
}

func TestParseHtml(t *testing.T) {
	data := bytes.NewBuffer([]byte(`
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Go 语言环境安装 | 菜鸟教程</title>
		<title>不要的</title>
		<title>多余的</title>
	</head>
</html>
`))
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		log.Fatal(err)
	}
	node := doc.Find("title").First().Text()
	t.Log(node)
}

func TestMatchInformation(t *testing.T) {
	// https://1257120875.vod2.myqcloud.com/0ef121cdvodtransgzp1257120875/3055695e5285890780828799271/v.f230.m3u8
	//u, err := url.Parse("https://1257120875.vod2.myqcloud.com/0ef121cdvodtransgzp1257120875/3055695e5285890780828799271/v.f230.m3u8")
	//u, err := url.Parse("https://1257120875.vod2.myqcloud.com/v.f230.mp4")
	u, err := url.Parse("https://1257120875.vod2.myqcloud.com/v/f230.m3u8")
	if err != nil {
		t.Fatal(err)
	}

	//parts := strings.Split(u.Path, "/")

	reg, err := regexp.Compile(`([^\/]+)(\.m3u8|\.mp4)$`)
	if err != nil {
		t.Fatal(err)
	}

	result := reg.FindString(u.Path)
	t.Log(len(result), result)

	html := `
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Go 语言环境安装 | 菜鸟教程</title>
		<title>不要的</title>
		<title>多余的</title>

f230.m3u8

	</head>
</html>

`

	regKey, _ := regexp.Compile(fmt.Sprintf("(%s)|(%s)|(%s)",
		"https://1257120875.vod2.myqcloud.com/v/f230.m3u8", "f230.m3u8", "f230"))
	t.Log(regKey.MatchString(html))

}
