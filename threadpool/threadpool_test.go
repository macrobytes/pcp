package threadpool

import (
	"container/list"
	"sync"
	"testing"
)

var buffer = list.New()
var mutex sync.Mutex

type testRunner struct {
	element int
}

func createRunner(i int) *testRunner {
	var runner testRunner
	runner.element = i
	return &runner
}

func (e *testRunner) Run() {
	mutex.Lock()
	buffer.PushBack(e.element)
	mutex.Unlock()
}

func TestExecutorLifeCycle(t *testing.T) {
	maxRoutines := 5
	executor := Init(maxRoutines)
	bufferLength := 100

	for i := 0; i < bufferLength; i++ {
		runner := createRunner(i)
		executor.Execute(runner)
		if executor.activeRoutines > maxRoutines {
			t.Errorf("Thread pool size of %v exceeded", maxRoutines)
		}
	}

	executor.Wait()
	if buffer.Len() != bufferLength {
		t.Errorf("Expected list size %v, got %v", bufferLength, buffer.Len())
	}
}
