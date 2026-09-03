package scheduler

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"chain_simulation/entities"

	"github.com/schollz/progressbar/v3"
)

type Scheduler struct {
	StartedAt      time.Time
	PendingEvents  []*entities.Event
	ExecutedEvents []*entities.Event
	ProgressBar    *progressbar.ProgressBar
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// Run executes events according to their StartTime relative to the beginning
// of this call. Events with the same StartTime retain their original order.
func Run(events []*entities.Event) error {
	return NewScheduler().Run(events)
}

func (scheduler *Scheduler) Run(events []*entities.Event) error {
	scheduler.StartedAt = time.Now()
	scheduler.PendingEvents = append([]*entities.Event(nil), events...)
	scheduler.ExecutedEvents = make([]*entities.Event, 0, len(events))
	scheduler.ProgressBar = newProgressBar(len(events), "remained-events")

	sort.SliceStable(scheduler.PendingEvents, func(left, right int) bool {
		return scheduler.PendingEvents[left].StartTime < scheduler.PendingEvents[right].StartTime
	})

	var executionErrors []error
	for len(scheduler.PendingEvents) > 0 {
		event := scheduler.PendingEvents[0]
		if waitDuration := event.StartTime - time.Since(scheduler.StartedAt); waitDuration > 0 {
			time.Sleep(waitDuration)
		}

		fmt.Printf("execute event %s\n", event.Action.String())
		if event.Handler == nil {
			executionErrors = append(executionErrors, fmt.Errorf("event %s has no handler", event.Action.String()))
		} else if err := event.Handler(); err != nil {
			executionErrors = append(executionErrors, fmt.Errorf("execute event %s: %w", event.Action.String(), err))
		}

		scheduler.PendingEvents = scheduler.PendingEvents[1:]
		scheduler.ExecutedEvents = append(scheduler.ExecutedEvents, event)
		_ = scheduler.ProgressBar.Add(1)
	}

	return errors.Join(executionErrors...)
}
