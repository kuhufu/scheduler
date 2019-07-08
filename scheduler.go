package scheduler

import (
	"sync"
	"time"
)

type Scheduler struct {
	start       chan struct{}
	stop        chan struct{}
	running     bool
	runningLock *sync.Mutex
}

func New() *Scheduler {
	return &Scheduler{
		start:       make(chan struct{}),
		stop:        make(chan struct{}),
		running:     false,
		runningLock: &sync.Mutex{},
	}
}

func (s *Scheduler) Stop() {
	s.runningLock.Lock()
	if !s.running {
		s.runningLock.Unlock()
		return
	}
	s.running = false
	s.runningLock.Unlock()
	close(s.stop)
}

func (s *Scheduler) Start() {
	s.runningLock.Lock()
	if s.running {
		s.runningLock.Unlock()
		return
	}
	s.running = true
	s.runningLock.Unlock()
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
