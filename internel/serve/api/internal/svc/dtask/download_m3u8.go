package dtask

import (
	"bytes"
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/svc/dtask/pool"
	"dv/internel/serve/api/internal/util/aes"
	"dv/internel/serve/api/internal/util/m3u8"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DownloadM3u8(pt *pool.Pool, mt model.Task, cfg config.Config) error {
	// todo 判断是否已经得到m3u8文件
	tLog := newTaskLog(pt.Ctx, mt)
	downCtrl := newDownload(pt.Ctx, tLog)
	req, err := parseTaskToHttp(mt)
	if err != nil {
		return err
	}

	fileM3u8Dir := filepath.Join(cfg.SaveDir, "m3u8")
	fileM3u8Name := mt.Name
	segments, aesKey, err := getM3u8ProtoFile(downCtrl, req, fileM3u8Dir, fileM3u8Name)
	if err != nil {
		return err
	}

	var playbackDuration float32 // 该视频总时间
	for _, segment := range segments {
		playbackDuration += segment.Duration
	}
	tLog.setScopeMax(int64(len(segments)))
	logx.Infof("%v 该电影时长 %v \n", mt.Name, m3u8.CalculationTime(playbackDuration))

	fileDir := filepath.Join(cfg.SaveDir, mt.Name)
	for index, segment := range segments {
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return err
		}
		link, err = link.Parse(segment.URI)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s_%05d", mt.Name, index)
		pathPart := strings.Split(link.Path, ".")
		if len(pathPart) > 0 {
			fileName = fmt.Sprintf("%s.%s", fileName, pathPart[len(pathPart)-1])
		}
		filePath := filepath.Join(fileDir, fileName)

		if err := pt.Submit(&pool.Cell{
			CellFunc: Control.listenExecute,
			Param: []interface{}{
				mt, pt, cfg, model.VideoTypeM3u8Single, filePath, aesKey,
			},
		}); err != nil {
			return err
		}

	}

	return nil
}

func getM3u8ProtoFile(downCtrl *download, req *http.Request, fileDir, fileName string) ([]*m3u8.Segment, []byte, error) {
	// 保存m3u8文件
	var saveM3u8SourceFile []byte
	// 密匙
	var aesKey []byte

	// 构建请求
	buf := bytes.NewBuffer(nil)
	if err := downCtrl.getSingle(defaultClient, req, buf); err != nil {
		return nil, nil, err
	}
	logx.Debug("m3u8_file ========> \n", buf.String())
	saveM3u8SourceFile = buf.Bytes()
	m3u8Data, err := m3u8.ParseM3u8Data(buf)
	if err != nil {
		return nil, nil, err
	}

	if len(m3u8Data.MasterPlaylist) != 0 {
		// 下载最高清的视频
		index := m3u8Data.GetMaxBandWidth()
		if index < 0 {
			return nil, nil, errors.New("解析失败")
		}
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, nil, err
		}
		link, err = link.Parse(m3u8Data.MasterPlaylist[index].URI)
		if err != nil {
			return nil, nil, err
		}

		request, err := http.NewRequest(http.MethodGet, link.String(), nil)
		if err != nil {
			return nil, nil, err
		}
		request.Header = req.Header
		return getM3u8ProtoFile(downCtrl, request, fileDir, fileName)
	}

	for _, key := range m3u8Data.Keys {
		if key.Method == m3u8.CryptMethodNONE {
			continue
		}
		// 获取加密密匙
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, nil, err
		}
		aesUrl, err := link.Parse(key.URI)
		if err != nil {
			return nil, nil, err
		}
		aesBuf := bytes.NewBuffer(nil)
		request, err := http.NewRequest(http.MethodGet, aesUrl.String(), nil)
		if err != nil {
			return nil, nil, err
		}
		request.Header = req.Header
		if err := downCtrl.getSingle(defaultClient, request, aesBuf); err != nil {
			return nil, nil, err
		}

		aesKey = aesBuf.Bytes()
		break
	}

	m3u8.SaveM3u8File(fileDir, fileName, saveM3u8SourceFile)
	return m3u8Data.Segments, aesKey, nil
}

func DownloadM3u8SingleVideo(params []interface{}) error {
	mt := params[0].(model.Task)
	pt := params[1].(*pool.Pool)
	//cfg := params[2].(config.Config)
	//videoType := params[3].(string)
	filePath := params[4].(string)
	aseKey := params[5].([]byte)

	buf := bytes.NewBuffer(nil)
	if err := downloadSingleWithWrite(buf, 0, DownloadChunkFlag, pt, mt); err != nil {
		return err
	}

	if len(aseKey) != 0 {
		data := aes.AESDecrypt(buf.Bytes(), aseKey)
		if data == nil {
			return errors.New("视频格式解析失败")
		}
		buf = bytes.NewBuffer(nil)
		buf.Write(data)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, buf); err != nil {
		return err
	}

	return nil
}
