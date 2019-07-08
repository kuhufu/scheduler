package scheduler

import (
	"time"
)

type Scheduler struct {
	start   chan struct{}
	stop    chan struct{}
	running bool
}

func New() *Scheduler {
	return &Scheduler{
		start:   make(chan struct{}),
		stop:    make(chan struct{}),
		running: false,
	}
}

func (s *Scheduler) Stop() {
	if !s.running {
		return
	}
	s.running = false
	close(s.stop)
}

func (s *Scheduler) Start() {
	if s.running {
		return
	}
	s.running = true
	close(s.start)
}

func (s *Scheduler) AddTimeoutFunc(duration time.Duration, cmd func()) {
	go func() {
		<-s.start
		timer := time.AfterFunc(duration, cmd)
		for {
			select {
			case <-s.stop:
				timer.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) AddIntervalFunc(duration time.Duration, cmd func()) {
	go func() {
		<-s.start
		tick := time.NewTicker(duration)
		tickC := tick.C
		for {
			select {
			case <-s.stop:
				tick.Stop()
				return
			case <-tickC:
				cmd()
			}
		}
	}()
}
