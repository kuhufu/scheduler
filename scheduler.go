package scheduler

import "time"

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
	s.stop <- struct{}{}
}

func (s *Scheduler) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.Run()
}

func (s *Scheduler) Run() {
	if s.running {
		return
	}
	s.running = true
	s.start <- struct{}{}
}

func (s *Scheduler) AddTimeoutFunc(duration time.Duration, cmd func()) {
	timer := time.AfterFunc(duration, cmd)
	go func() {
		for{
			select {
			case <-s.stop:
				timer.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) AddIntervalFunc(duration time.Duration, cmd func()) {
	tick := time.NewTicker(duration)
	tickC := tick.C
	//tick := time.Tick(duration)
	go func() {
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
