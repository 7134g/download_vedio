package dtask

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	DownloadSingleFlag = 1 // 下载整个完整文件
	DownloadChunkFlag  = 2 // 下载整个文件块
)

type download struct {
	status bool // 状态

	tLog *taskLog
}

func newDownload(ctx context.Context, tLog *taskLog) *download {
	d := &download{
		status: false,
		tLog:   tLog,
	}
	stopChan := make(chan struct{})
	defer close(stopChan)
	go d.stop(ctx, stopChan)

	return d
}

func (d *download) stop(ctx context.Context, stopChan chan struct{}) {
	select {
	case <-ctx.Done():
		d.status = false
	case <-stopChan:
		return
	}
}

func (d *download) get(client *http.Client, req *http.Request, write io.Writer) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("响应为 %s", resp.Status))
	}

	if resp.ContentLength == 0 {
		logx.Error(fmt.Sprintf("%s 跳过数据内容大于等于文件大小，因此不下载\n", d.tLog.taskName))
		return nil, nil
	}

	return resp, nil
}

// 读取整个
func (d *download) getSingle(client *http.Client, req *http.Request, write io.Writer) error {
	resp, err := d.get(client, req, write)
	if resp != nil {
		defer resp.Body.Close()
	} else {
		return nil
	}
	if err != nil {
		return err
	}

	ctxRange := resp.Header.Get("Content-Range")
	if len(ctxRange) == 0 {
		d.tLog.setScopeMax(resp.ContentLength)
	} else {
		begin := strings.Index(ctxRange, " ")
		end := strings.Index(ctxRange, "-")
		haveLengthString := ctxRange[begin+1 : end]
		haveLength, err := strconv.Atoi(haveLengthString)
		if err != nil {
			return err
		}
		completeFileSizeString := ctxRange[strings.LastIndex(ctxRange, "/")+1:]
		completeFileSize, err := strconv.Atoi(completeFileSizeString)
		if err != nil {
			return err
		}

		d.tLog.setScopeMax(int64(completeFileSize))
		d.tLog.setScope(int64(haveLength), int64(haveLength))
	}

	return d.rwSingle(resp.Body, write)
}

func (d *download) rwSingle(read io.Reader, write io.Writer) error {
	bs := make([]byte, 1048576) // 每次读取http内容的大小(1mb)

	for {
		if d.status {
			return nil
		}

		rn, err := read.Read(bs)
		if err != nil {
			if err == io.EOF {
				// 完成
				_, _ = write.Write(bs[:rn])
				d.tLog.incScope(int64(rn), int64(rn))
				return nil
			}
			return err
		}

		_, err = write.Write(bs[:rn])
		if err != nil {
			return err
		}
		d.tLog.incScope(int64(rn), int64(rn))
	}
}

// 读取块
func (d *download) getChunk(client *http.Client, req *http.Request, write io.Writer) error {
	resp, err := d.get(client, req, write)
	if resp != nil {
		defer resp.Body.Close()
	} else {
		d.tLog.incScope(1, 0)
		return nil
	}
	if err != nil {
		return err
	}
	return d.rwChunk(resp.Body, write)
}

func (d *download) rwChunk(read io.Reader, write io.Writer) error {
	bs := make([]byte, 1048576) // 每次读取http内容的大小(1mb)

	for {
		if d.status {
			return nil
		}
		rn, err := read.Read(bs)
		if err != nil {
			if err == io.EOF {
				// 完成
				_, _ = write.Write(bs[:rn])
				d.tLog.incScope(1, int64(rn))
				return nil
			}
			return err
		}

		_, err = write.Write(bs[:rn])
		if err != nil {
			return err
		}
		d.tLog.incScope(0, int64(rn))
	}
}
