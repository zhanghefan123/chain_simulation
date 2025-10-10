package thread_manager

import "sync"

var ThreadManagerInstance = NewThreadManager()

type ThreadManager struct {
	WaitGroup sync.WaitGroup
}

func NewThreadManager() *ThreadManager {
	return &ThreadManager{
		WaitGroup: sync.WaitGroup{},
	}
}

func (tm *ThreadManager) Add() {
	tm.WaitGroup.Add(1)
}

func (tm *ThreadManager) Wait() {
	tm.WaitGroup.Wait()
}

func (tm *ThreadManager) Done() {
	tm.WaitGroup.Done()
}
