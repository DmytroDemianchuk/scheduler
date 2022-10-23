package scheduler

import (
	"context"
	"sync"
	"time"
)

type Job func(ctx context.Context)

type Scheduler struct {
	wg           *sync.WaitGroup
	cancellation []context.CancelFunc
}

func NewScheduler() *Scheduler {
	return *Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
	}
}

func (s *Scheduler) Add(ctx context.Context, j Job, interval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellation = append(s.cancellation, cancel)

	s.wg.Add(1)
	go s.process(ctx, j, interval)
}

func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellation {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) precess(ctx context.Context, j Job, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			j(ctx)
		case <-ctx.Done():
			s.wg.Done()
			ticker.Stop()
			return
		}
	}
}
