package threadpool

import (
	"sync"
)

// Executor - manages and executes go routines
type Executor struct {
	waitGroup      sync.WaitGroup
	mutex          sync.Mutex
	activeRoutines int
	maxRoutines    int
}

func (e *Executor) done() {
	e.activeRoutines--
	e.waitGroup.Done()
}

func (e *Executor) execute(task Runnable) {
	task.Run()
	e.done()
}

// Execute - Runs the specified function as a go-routine when space is
// available in pool.
func (e *Executor) Execute(task Runnable) {
	e.waitGroup.Add(1)
	for {
		e.mutex.Lock()
		if e.activeRoutines < e.maxRoutines {
			e.activeRoutines++
			go e.execute(task)
			e.mutex.Unlock()
			break
		}
		e.mutex.Unlock()
	}
}

// Wait - Waits for all executed go-routines to complete.
func (e *Executor) Wait() {
	e.waitGroup.Wait()
}

// Init - initializes the Executor
func Init(maxRoutines int) *Executor {
	var executor Executor
	executor.maxRoutines = maxRoutines
	executor.activeRoutines = 0
	return &executor
}
