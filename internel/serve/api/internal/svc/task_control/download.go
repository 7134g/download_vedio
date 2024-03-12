package task_control

import (
	"bytes"
	"dv/internel/serve/api/internal/util/m3u8"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type download struct {
	req *http.Request

	t        model.Task
	fileDir  string // 目录
	fileName string // 文件名

	//fileSize      int64         // 目前文件大小
	totalFileSize int64         // 文件总大小
	stop          chan struct{} // 打印进度
}

func newDownload(t model.Task, fileDir, fileName string, ext bool) *download {
	_ = os.MkdirAll(fileDir, 0700)
	var fn string
	if !ext {
		fn = fmt.Sprintf("%s.mp4", fileName)
	} else {
		fn = fileName
	}

	return &download{
		t: t,

		fileDir:  fileDir,
		fileName: fn,
		stop:     make(chan struct{}),
	}
}

func buildKey(id uint, name string, childFlag ...string) string {
	key := fmt.Sprintf("%d_%s", id, name)
	if len(childFlag) > 0 {
		for _, child := range childFlag {
			key = fmt.Sprintf("%s_%s", key, child)
		}
	}

	return key
}

func (d *download) getM3u8File(client *http.Client, req *http.Request) ([]*m3u8.Segment, error) {
	// 保存m3u8文件
	var saveM3u8SourceFile []byte

	// 构建请求
	buf := bytes.NewBuffer(nil)
	if err := d.get(client, req, buf); err != nil {
		return nil, err
	}
	logx.Debug(buf.String())
	saveM3u8SourceFile = buf.Bytes()
	m3u8Data, err := m3u8.ParseM3u8Data(buf)
	if err != nil {
		return nil, err
	}

	if len(m3u8Data.MasterPlaylist) != 0 {
		// 下载最高清的视频
		index := m3u8Data.GetMaxBandWidth()
		if index < 0 {
			return nil, errors.New("解析失败")
		}
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, err
		}
		link, err = link.Parse(m3u8Data.MasterPlaylist[index].URI)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodGet, link.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header = req.Header
		return d.getM3u8File(client, request)
	}

	for _, key := range m3u8Data.Keys {
		if key.Method == m3u8.CryptMethodNONE {
			continue
		}
		// 获取加密密匙
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, err
		}
		aesUrl, err := link.Parse(key.URI)
		if err != nil {
			return nil, err
		}
		aesBuf := bytes.NewBuffer(nil)
		request, err := http.NewRequest(http.MethodGet, aesUrl.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header = req.Header
		if err := d.get(client, request, aesBuf); err != nil {
			return nil, err
		}

		table.CryptoVideoTable.Set(d.t.ID, aesBuf.Bytes())
		break
	}

	m3u8.SaveM3u8File(d.fileDir, d.t.Name, saveM3u8SourceFile)
	return m3u8Data.Segments, nil
}

func (d *download) get(client *http.Client, req *http.Request, write io.Writer) error {
	d.stop = make(chan struct{})
	// 构建请求
	logx.Debug(d.t.Name, req.URL.String())
	//x, _ := json.Marshal(req.Header)
	//logx.Debug(string(x))
	req.Header = http.Header{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s", resp.Status))
	}

	if resp.ContentLength == 0 {
		return errors.New(fmt.Sprintf("%s 跳过数据内容大于等于文件大小，因此不下载\n", d.fileName))
	}

	ctxRange := resp.Header.Get("Content-Range")
	if len(ctxRange) == 0 {
		// 记录文件总大小
		d.recordDownloadMaxLength(uint(resp.ContentLength))
	} else {
		// Content-Range: bytes 1629222-5510871/5510872 取 5510872
		// 5510872 指的是文件总大小

		begin := strings.Index(ctxRange, " ")
		end := strings.Index(ctxRange, "-")
		haveLengthString := ctxRange[begin+1 : end]
		haveLength, err := strconv.Atoi(haveLengthString)
		if err != nil {
			return err
		}
		d.recordDownloadNow(uint(haveLength))

		completeFileSizeString := ctxRange[strings.LastIndex(ctxRange, "/")+1:]
		completeFileSize, err := strconv.Atoi(completeFileSizeString)
		if err != nil {
			return err
		}
		d.recordDownloadMaxLength(uint(completeFileSize))
	}

	return d.rw(resp.Body, write)
}

func (d *download) rw(read io.Reader, write io.Writer) error {
	defer close(d.stop)

	bs := make([]byte, 1048576) // 每次读取http内容的大小(1mb)
	for {
		rn, err := read.Read(bs)
		if err != nil {
			if err == io.EOF {
				// 完成
				//d.fileSize += int64(rn)
				d.recordDownloadNow(uint(rn))
				_, _ = write.Write(bs[:rn])
				return nil
			}
			return err
		}

		d.recordDownloadNow(uint(rn))

		//d.fileSize += int64(rn)
		_, err = write.Write(bs[:rn])
		if err != nil {
			return err
		}
	}
}

func (d *download) recordDownloadNow(value uint) {
	if d.t.VideoType != model.VideoTypeM3u8 {
		table.DownloadTaskByteLength.Inc(d.t.ID, value)
		table.DownloadTimeSince.Set(d.t.ID, value)
	}
}

func (d *download) recordDownloadMaxLength(value uint) {
	if d.t.VideoType != model.VideoTypeM3u8 {
		table.DownloadTaskMaxLength.Set(d.t.ID, value)
	}
}
