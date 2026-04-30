package scheduler

import (
	"chain_simulation/entities"
	"chain_simulation/modules/progress_bar"
	"chain_simulation/modules/thread_manager"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"sync"
	"time"
)

var SchedulerInstance = NewScheduler()

type Scheduler struct {
	CurrentTime       time.Time
	StopQueue         chan struct{}
	EventList         []*entities.Event
	ExecutedEventList []*entities.Event
	WaitGroup         *sync.WaitGroup
	ProgressBar       *progressbar.ProgressBar
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		StopQueue:         make(chan struct{}),
		EventList:         make([]*entities.Event, 0),
		ExecutedEventList: make([]*entities.Event, 0),
		WaitGroup:         &sync.WaitGroup{},
		ProgressBar:       nil,
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

func SetEventsIntoScheduler(events []*entities.Event) {
	SchedulerInstance.EventList = events
	SchedulerInstance.ProgressBar = progress_bar.NewProgressBar(len(SchedulerInstance.EventList), "remained-events")
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
	var ticker = time.NewTicker(time.Millisecond * 1000)
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

			// 找到第一个 event
			if len(eventsToRemain) > 0 {
				fmt.Printf("first event %s time remained = %v\n", eventsToRemain[0].Action.String(), eventsToRemain[0].StartTime-time.Since(s.CurrentTime))
			}

			s.EventList = eventsToRemain

			// 2. 进行执行
			if len(eventsToExecute) > 0 {
				for _, event := range eventsToExecute {
					// 每个事件进行单独的执行
					fmt.Printf("execute event\n")
					err := event.Handler()
					if err != nil {
						fmt.Printf("Error executing event: %v\n", err)
					}
					s.ExecutedEventList = append(s.ExecutedEventList, event)

					// 进行进度条的更新
					_ = s.ProgressBar.Add(1)
				}
			}
		}
	}
}

func (s *Scheduler) Stop() {
	s.StopQueue <- struct{}{}
	s.WaitGroup.Wait()
}
