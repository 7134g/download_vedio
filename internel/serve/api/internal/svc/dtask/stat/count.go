package stat

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

var StatisticsControl Statistics

type Statistics struct {
	lock sync.Mutex
	cfg  config.Config

	errorCount map[int]int // 任务错误次数
}

func Init(cfg config.Config) {
	StatisticsControl = Statistics{
		lock:       sync.Mutex{},
		cfg:        cfg,
		errorCount: map[int]int{},
	}
}

func (s *Statistics) IncError(t model.Task, err error) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	logx.Errorw("错误",
		logx.Field("err", err),
		logx.Field("name", t.Name),
	)
	if errorCount, exist := s.errorCount[t.ID]; exist {
		if errorCount >= s.cfg.TaskErrorMaxCount {
			return false
		} else {
			s.errorCount[t.ID]++
			return true
		}
	} else {
		s.errorCount[t.ID] = 1
		return true
	}
}

func (s *Statistics) DelErrorCount(taskId int) {
	delete(s.errorCount, taskId)
}
