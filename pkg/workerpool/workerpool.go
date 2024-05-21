package workerpool

import "sync"

type Task func()

type WorkerPool struct {
	tasks     chan Task
	jobCount  int
	wg        sync.WaitGroup
	addTaskWg sync.WaitGroup
}

func New(jobCount int) *WorkerPool {
	return &WorkerPool{
		jobCount: jobCount,
		tasks:    make(chan Task),
	}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.jobCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.addTaskWg.Add(1)
	go func() {
		defer wp.addTaskWg.Done()
		wp.tasks <- task
	}()
}

func (wp *WorkerPool) Close() {
	wp.addTaskWg.Wait()
	close(wp.tasks)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		task()
	}
}
