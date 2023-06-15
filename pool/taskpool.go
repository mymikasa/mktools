package pool

import (
	"sync"
	"time"
)

const (
	idleTimeout = 2 * time.Second
)

type TaskPool struct {
	maxWorkers  int
	taskQueue   chan func()
	workerQueue chan func()
	stoppedChan chan struct{}
	stopSignal  chan struct{}
	waitQueue   chan struct{}
	stoped      bool
	stopLock    sync.Mutex
	stopOnce    sync.Once

	wait bool
}

// Size 返回当前TaskPool能开启的最大数量
func (p *TaskPool) Size() int {
	return p.maxWorkers
}

func (p *TaskPool) Stop() {

}

func (p *TaskPool) stop(wait bool) {
	p.stopOnce.Do(func() {
		close(p.stopSignal)

		p.stopLock.Lock()
		p.stoped = true
		p.stopLock.Unlock()
		p.wait = wait
		close(p.taskQueue)
	})
	<-p.stoppedChan
}
