package workerpool

import "sync"

type Task func()

type WorkerPool struct {
	tasks chan Task
	wg    sync.WaitGroup
}

const workerPoolQueueSize = 5

func New(workerCount int) *WorkerPool {
	return NewWithQS(workerCount, workerPoolQueueSize)
}

func NewWithQS(workerCount int, queueCount int) *WorkerPool {
	wp := &WorkerPool{
		tasks: make(chan Task, queueCount),
	}
	for range workerCount {
		wp.wg.Add(1)
		go wp.worker()
	}
	return wp
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Close() {
	close(wp.tasks)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		task()
	}
}
