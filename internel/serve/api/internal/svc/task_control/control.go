package task_control

import (
	"context"
	"dv/internel/serve/api/internal/util/calc"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"log"
	"path/filepath"
	"time"
)

func (c *TaskControl) GetStatus() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.running
}

func (c *TaskControl) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cancel()
	close(c.printStop)
	close(c.vacancy)
	c.vacancy = make(chan struct{}, len(c.vacancy))
	c.running = false
}

func (c *TaskControl) start() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.running = true
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.printStop = make(chan struct{})
}

func (c *TaskControl) incDoneCount(d *download) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.doneCount++
	table.DownloadTaskMaxLength.Set(d.t.ID, 0)
	table.DownloadTaskByteLength.Set(d.t.ID, 0)
	table.DownloadTimeSince.Set(d.t.ID, 0)
}

// 所有任务进度
func (c *TaskControl) printDownloadProgress(t model.Task) {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-c.printStop:
			return
		case <-ticker.C:
			maxLength, exist := table.DownloadTaskMaxLength.Get(t.ID)
			if !exist {
				continue
			}
			nowDownloadDataLen, exist := table.DownloadTaskByteLength.Get(t.ID)
			if !exist {
				continue
			}

			downloadTimeSince, _ := table.DownloadTimeSince.Get(t.ID)
			speed, unit := calc.UnitReturn(float64(downloadTimeSince))
			score := float64(nowDownloadDataLen) / float64(maxLength) * 100
			log.Println(fmt.Sprintf("%s 进度：%.2f 速度：%.2f %s/s",
				t.Name,
				score,
				speed/3, unit,
			))
		}
	}
}

func (c *TaskControl) submit(fn particleFunc, params []any) {
	c.wg.Add(1)
	select {
	case c.vacancy <- struct{}{}:
	case <-c.ctx.Done():
		logx.Info("cancel stop")
		return
	}

	d := params[0].(*download)
	go threading.GoSafe(func() {
		defer c.submitDone()
		if d == nil {
			return
		}

		//keyPart := strings.Split(d.key, "_")
		//taskId, _ := strconv.Atoi(keyPart[0])

		key := buildKey(d.t.ID, d.t.Name)
		if err := fn([]any{d}); err != nil {
			table.IncErrCount(key)
			if table.GetErrCount(key) >= tcConfig.TaskErrorMaxCount {
				_ = tasKModel.UpdateStatus(d.t.ID, model.StatusError)
				logx.Error(d.t.Name, "任务失败", err)
			} else {
				logx.Errorw(
					"control error message",
					logx.Field("retry_count", table.GetErrCount(key)),
					logx.Field("key", key),
					logx.Field("error", err),
					logx.Field("sleep_time", tcConfig.TaskErrorDuration),
				)
				time.Sleep(time.Second * time.Duration(tcConfig.TaskErrorDuration))
				go c.submit(fn, []any{d}) // 重试
			}
			return
		}

		c.incDoneCount(d)
		if d.t.VideoType == model.VideoTypeMp4 {
			_ = tasKModel.UpdateStatus(d.t.ID, model.StatusSuccess)
			logx.Info(d.t.Name, " is done")
		}

		return
	})
}

func (c *TaskControl) submitDone() {
	c.wg.Done()
	<-c.vacancy
}

func (c *TaskControl) Run(tasks []model.Task) {
	logx.Info("running ......")
	logx.Info(filepath.Abs(tcConfig.SaveDir))
	defer c.Stop()
	c.start()

	//go c.printDownloadProgress(uint(len(tasks)))

	for _, m := range tasks {
		go c.printDownloadProgress(m)
		w := newWork(m)
		particle, d := w.parseTask()
		if particle == nil {
			continue
		}

		_ = tasKModel.UpdateStatus(m.ID, model.StatusRunning)
		logx.Infof("=========> 任务开始：%s", w.task.Name)
		c.submit(particle, []any{d})
	}
	c.wg.Wait()

	logx.Info("所有任务已经结束 <=========")
}
