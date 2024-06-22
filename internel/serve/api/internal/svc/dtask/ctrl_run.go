package dtask

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/dao"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/svc/dtask/pool"
	"time"
)

func (c *control) running(t model.Task, cfg config.Config) error {
	if err := dao.TaskDao.UpdateStatus(t.ID, model.StatusRunning); err != nil {
		return err
	}

	var size int32 = 1
	switch t.Type {
	case model.VideoTypeM3u8:
		size = cfg.TaskM3u8Concurrency + 1
	}

	p, cancel := pool.NewPool(size, false, time.Second)
	c.task[t.ID] = t
	c.pools[t.ID] = p
	c.cancels[t.ID] = cancel

	if err := p.Submit(&pool.Cell{
		CellFunc: c.listenExecute,
		Param: []interface{}{
			t, p, cfg, t.Type,
		},
	}); err != nil {
		return err
	}

	return nil
}
