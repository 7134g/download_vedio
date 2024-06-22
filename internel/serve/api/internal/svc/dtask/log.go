package dtask

import "C"
import (
	"context"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/util/encoding"
	"fmt"
	"sync"
	"time"
)

var ControlLog controlLog

type controlLog struct {
	lock   sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	TaskLogMap map[int]*taskLog

	cacheChan chan string
	CacheLog  []string
}

func init() {
	ControlLog.Init()
}

func (c *controlLog) Init() {
	if c.cancel != nil {
		c.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	ControlLog = controlLog{
		lock:       sync.Mutex{},
		ctx:        ctx,
		cancel:     cancel,
		TaskLogMap: make(map[int]*taskLog),
		cacheChan:  make(chan string, 10),
		CacheLog:   make([]string, 0),
	}

	go ControlLog.Listen()
}

func (c *controlLog) Listen() {
	for {
		select {
		case message := <-c.cacheChan:
			c.CacheLog = append(c.CacheLog, message)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *controlLog) Close() {
	c.cancel()
}

func (c *controlLog) GetTaskLog(taskId int) *taskLog {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, exist := c.TaskLogMap[taskId]
	if exist {
		return value
	}
	return nil
}

func (c *controlLog) AddTaskLog(l *taskLog) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.TaskLogMap[l.taskId] = l
}

func (c *controlLog) DelTaskLog(taskId int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.TaskLogMap, taskId)
}

type taskLog struct {
	taskId   int    // 任务id
	taskName string // 任务名称

	loading         int64 // 下载进度 map[taskid]scope
	lastLoading     int64 // 上次记录的进度，用于计算速度 map[taskid]scope
	chunkLength     int64 // 上次下载的区块大小
	lastChunkLength int64 // 本次下载的区块大小
	target          int64 // 目标进度 map[taskid]scopeMax

	intervalPrintTime int // 打印间隔时间
}

func newTaskLog(ctx context.Context, t model.Task) *taskLog {
	l := ControlLog.GetTaskLog(t.ID)
	if l != nil {
		return l
	}
	l = &taskLog{
		taskId:            t.ID,
		taskName:          t.Name,
		loading:           0,
		lastLoading:       0,
		chunkLength:       0,
		lastChunkLength:   0,
		target:            0,
		intervalPrintTime: 3,
	}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(l.intervalPrintTime))
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.print()
		}
	}()

	ControlLog.AddTaskLog(l)
	return l
}

func (c *taskLog) GetLoading() float64 {
	return float64(c.loading) / float64(c.target)
}

func (c *taskLog) setScopeMax(scopeMax int64) {
	c.target = scopeMax
}

func (c *taskLog) setScope(scope int64, speed int64) {
	c.lastLoading = scope
	c.loading = scope
	c.lastChunkLength = speed
	c.chunkLength = speed
}

func (c *taskLog) incScope(scope, speed int64) {
	c.lastLoading = c.loading
	c.loading = c.loading + scope
	c.lastChunkLength = c.chunkLength
	c.chunkLength = c.chunkLength + speed
	return
}

func (c *taskLog) print() {
	if c.loading == 0 {
		return
	}

	loading := c.loading
	lastLoading := c.lastLoading
	target := c.target
	name := c.taskName
	nowChunk := c.chunkLength
	lastChunk := c.lastChunkLength

	if loading == lastLoading {
		// 进度无变化
		return
	}
	if loading >= target {
		// 完成
		return
	}

	load := float64(loading) / float64(target)
	interval := nowChunk - lastChunk
	speed := encoding.ByteCountIEC(int64(float64(interval) / float64(c.intervalPrintTime)))
	message := fmt.Sprintf("=======> 任务：%s 当前下载进度：%s , 速度 %s, ID: %d",
		name,
		fmt.Sprintf("%2.2f", load)+" %/ 100 %",
		speed,
		c.taskId,
	)
	ControlLog.CacheLog = append(ControlLog.CacheLog, message)
	fmt.Println(message)
}
