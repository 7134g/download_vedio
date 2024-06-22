package dtask

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/task/pool"
	"os"
	"path/filepath"
)

func DownloadMp4(pt *pool.Pool, t model.Task, cfg config.Config) error {
	savePath := filepath.Join(cfg.SaveDir, t.Name)
	var flag = os.O_RDWR | os.O_CREATE | os.O_APPEND

	return DownloadSingleVideo(savePath, flag, pt, t)
}
