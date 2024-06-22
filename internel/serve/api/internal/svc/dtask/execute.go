package dtask

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/dao"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/task/pool"
	"dv/internel/serve/api/internal/task/stat"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (c *control) listenExecute(params []interface{}) {
	t := params[0].(model.Task)

	err := c.execute(params)
	if err != nil {
		if status := stat.StatisticsControl.IncError(t, err); status {
			time.Sleep(time.Duration(c.cfg.TaskErrorDuration) * time.Second)
			c.listenExecute(params)
		} else {
			// 达到了最大错误值
			logx.Error("error_max_count",
				logx.Field("id", t.ID),
				logx.Field("err", err),
				logx.Field("db_err", dao.TaskDao.UpdateStatus(t.ID, model.StatusError)),
			)
		}
		return
	}

	c.Done(t.ID)
}

func (c *control) execute(params []interface{}) error {
	t := params[0].(model.Task)
	pt := params[1].(*pool.Pool)
	cfg := params[2].(config.Config)
	videoType := params[3].(string)

	switch videoType {
	case model.VideoTypeMp4:
		return DownloadMp4(pt, t, cfg)
	case model.VideoTypeM3u8:
		return DownloadM3u8(pt, t, cfg)
	case model.VideoTypeM3u8Single:
		return DownloadM3u8SingleVideo(params)
	default:
		return errors.New("type error")
	}

	return nil
}
