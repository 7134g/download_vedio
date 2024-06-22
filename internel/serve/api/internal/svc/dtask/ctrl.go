package dtask

import (
	"context"
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/dao"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/svc/dtask/pool"
	"dv/internel/serve/api/internal/svc/dtask/stat"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

// 管理任务所有开始和停止

var Control control

type control struct {
	lock sync.Mutex
	cfg  config.Config

	status       bool  // 是否已经运行
	runningCount int32 // 正在执行的任务

	task    map[int]model.Task         // 任务表
	pools   map[int]*pool.Pool         // 工作表
	cancels map[int]context.CancelFunc // 取消
}

func (c *control) Init(cfg config.Config) {
	c.cfg = cfg

	c.status = false
	c.runningCount = 0

	c.task = map[int]model.Task{}
	c.pools = map[int]*pool.Pool{}
	c.cancels = map[int]context.CancelFunc{}

	// 初始化需要的
	InitDefaultHttp(cfg)
}

func (c *control) CheckStatus() bool {
	if !c.status {
		return false
	}

	if c.runningCount == 0 {
		c.status = false
		return false
	}

	return true
}

func (c *control) SetCfg(cfg config.Config) {
	c.cfg = cfg
}

func (c *control) Start(tasks []model.Task) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.status {
		return errors.New("running")
	}

	c.status = true
	for _, t := range tasks {
		if t.Status == model.StatusSuccess {
			continue
		}

		if err := dao.TaskDao.UpdateStatus(t.ID, model.StatusRunning); err != nil {
			logx.Error(err, logx.Field("id", t.ID))
			continue
		}
		go func(t model.Task) {
			for {
				c.runningCount++
				if c.runningCount >= c.cfg.TaskConcurrency {
					time.Sleep(time.Second * 5)
					continue
				}
				err := c.running(t, c.cfg)
				if err == nil {
					return
				}
				logx.Error(err, logx.Field("id", t.ID))
			}
		}(t)
	}

	return nil
}

func (c *control) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cancelFunc := range c.cancels {
		cancelFunc()
	}

	for _, p := range c.pools {
		p.Close()
	}

	for _, task := range c.task {
		err := dao.TaskDao.UpdateStatus(task.ID, model.StatusStop)
		if err != nil {
			logx.Error(err)
		}
	}

	c.Init(c.cfg)
}

func (c *control) Add(d model.Task, cfg config.Config) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exist := c.task[d.ID]; exist {
		return errors.New(fmt.Sprintf("exist %d", d.ID))
	}

	if _, exist := c.pools[d.ID]; exist {
		return errors.New(fmt.Sprintf("exist %d", d.ID))
	}

	return c.running(d, cfg)
}

func (c *control) Done(id int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.runningCount--
	_ = dao.TaskDao.Delete(id)
	c.pools[id].Close()
	c.cancels[id]()

	delete(c.task, id)
	delete(c.pools, id)
	delete(c.cancels, id)
	stat.StatisticsControl.DelErrorCount(id)
}
