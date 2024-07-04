package logging

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

// HookExecuter 接口定义了日志执行器的方法
type HookExecuter interface {
	Exec(extra map[string]string, b []byte) error
	Close() error
}

// hookOptions 结构体定义了 Hook 的选项
type hookOptions struct {
	maxJobs    int
	maxWorkers int
	extra      map[string]string
}

// HookOption 是一个用于设置 Hook 参数的函数类型
type HookOption func(*hookOptions)

// SetHookMaxJobs 设置最大任务数
func SetHookMaxJobs(maxJobs int) HookOption {
	return func(o *hookOptions) {
		o.maxJobs = maxJobs
	}
}

// SetHookMaxWorkers 设置最大工作线程数
func SetHookMaxWorkers(maxWorkers int) HookOption {
	return func(o *hookOptions) {
		o.maxWorkers = maxWorkers
	}
}

// SetHookExtra 设置额外的日志信息
func SetHookExtra(extra map[string]string) HookOption {
	return func(o *hookOptions) {
		o.extra = extra
	}
}

// NewHook 创建一个新的 Hook 实例
func NewHook(exec HookExecuter, opt ...HookOption) *Hook {
	opts := &hookOptions{
		maxJobs:    1024,
		maxWorkers: 2,
	}

	for _, o := range opt {
		o(opts)
	}

	wg := new(sync.WaitGroup)
	wg.Add(opts.maxWorkers)

	h := &Hook{
		opts: opts,
		q:    make(chan []byte, opts.maxJobs),
		wg:   wg,
		e:    exec,
	}
	h.dispatch()
	return h
}

// Hook 结构体定义了一个日志钩子，用于将日志发送到 MongoDB
type Hook struct {
	opts   *hookOptions
	q      chan []byte
	wg     *sync.WaitGroup
	e      HookExecuter
	closed int32
}

// dispatch 启动工作线程处理日志
func (h *Hook) dispatch() {
	for i := 0; i < h.opts.maxWorkers; i++ {
		go func() {
			defer func() {
				h.wg.Done()
				if r := recover(); r != nil {
					log.Println("Recovered from panic in logger hook:", r)
				}
			}()

			for data := range h.q {
				err := h.e.Exec(h.opts.extra, data)
				if err != nil {
					fmt.Println("Failed to write entry:", err.Error())
				}
			}
		}()
	}
}

// Write 将日志写入队列
func (h *Hook) Write(p []byte) (int, error) {
	if atomic.LoadInt32(&h.closed) == 1 {
		return len(p), nil
	}

	if len(h.q) == h.opts.maxJobs {
		log.Println("Too many jobs, waiting for queue to be empty, discard")
		return len(p), nil
	}

	data := make([]byte, len(p))
	copy(data, p)
	h.q <- data

	return len(p), nil
}

// Flush 关闭队列并等待所有工作线程完成
func (h *Hook) Flush() {
	if atomic.LoadInt32(&h.closed) == 1 {
		return
	}
	atomic.StoreInt32(&h.closed, 1)
	close(h.q)
	h.wg.Wait()
	err := h.e.Close()
	if err != nil {
		fmt.Println("Failed to close logger hook:", err.Error())
	}
}
