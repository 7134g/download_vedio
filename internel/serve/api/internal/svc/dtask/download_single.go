package dtask

import (
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/svc/dtask/pool"
	"fmt"
	"io"
	"os"
)

func DownloadSingleVideo(savePath string, openFlag int, pt *pool.Pool, mt model.Task) error {
	videoFile, err := os.OpenFile(savePath, openFlag, os.ModePerm)
	if err != nil {
		return err
	}
	info, err := videoFile.Stat()
	if err != nil {
		return err
	}

	return downloadSingleWithWrite(videoFile, info.Size(), DownloadSingleFlag, pt, mt)
}

func downloadSingleWithWrite(write io.Writer, fileSize int64, df int, pt *pool.Pool, t model.Task) error {
	tLog := newTaskLog(pt.Ctx, t)
	downCtrl := newDownload(pt.Ctx, tLog)
	req, err := parseTaskToHttp(t)
	if err != nil {
		return err
	}
	if fileSize > 0 {
		// 断点续传
		tLog.setScope(fileSize, 0)
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fileSize))
	}

	switch df {
	case DownloadSingleFlag:
		if err := downCtrl.getSingle(defaultClient, req, write); err != nil {
			return err
		}
	case DownloadChunkFlag:
		if err := downCtrl.getChunk(defaultClient, req, write); err != nil {
			return err
		}
	}

	return nil
}
