package scheduler

import (
	"chain_simulation/entities"
	"chain_simulation/modules/thread_manager"
	"fmt"
	"sync"
	"time"
)

var SchedulerInstance = NewScheduler()

type Scheduler struct {
	CurrentTime time.Time
	StopQueue   chan struct{}
	EventList   []*entities.Event
	WaitGroup   *sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		StopQueue: make(chan struct{}),
		EventList: make([]*entities.Event, 0),
		WaitGroup: &sync.WaitGroup{},
	}
}

func StartScheduler() {
	thread_manager.ThreadManagerInstance.Add()
	go func() {
		defer func() {
			thread_manager.ThreadManagerInstance.Done()
		}()
		SchedulerInstance.Start()
	}()
}

func AddEventIntoScheduler(event *entities.Event) {
	SchedulerInstance.AddEvent(event)
}

func SetEventsIntoScheduler(events []*entities.Event) {
	SchedulerInstance.EventList = events
}

func StopScheduler() {
	SchedulerInstance.StopQueue <- struct{}{}
	SchedulerInstance.WaitGroup.Wait()
}

func (s *Scheduler) AddEvent(event *entities.Event) {
	s.EventList = append(s.EventList, event)
}

func (s *Scheduler) Start() {
	s.WaitGroup.Add(1)
	defer func() {
		s.WaitGroup.Done()
	}()
	s.CurrentTime = time.Now()
	var ticker = time.NewTicker(time.Millisecond * 100)
ForLoop:
	for {
		select {
		case <-s.StopQueue:
			break ForLoop
		case <-ticker.C:
			// 1. 计算要进行执行和保留的 event
			eventsToExecute := make([]*entities.Event, 0)
			eventsToRemain := make([]*entities.Event, 0)
			for _, event := range s.EventList {
				if time.Since(s.CurrentTime) > event.StartTime {
					eventsToExecute = append(eventsToExecute, event)
				} else {
					eventsToRemain = append(eventsToRemain, event)
				}
			}
			s.EventList = eventsToRemain
			// 2. 进行执行
			if len(eventsToExecute) > 0 {
				for _, event := range eventsToExecute {
					err := event.Handler()
					if err != nil {
						fmt.Printf("Error executing event: %v\n", err)
					}
				}
			}
		}
	}
}

func (s *Scheduler) Stop() {
	s.StopQueue <- struct{}{}
	s.WaitGroup.Wait()
}
