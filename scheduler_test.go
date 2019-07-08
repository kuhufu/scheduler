package scheduler

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

var s = New()

func TestScheduler_AddIntervalFunc(t *testing.T) {

	s.AddIntervalFunc(time.Second, func() {
		log.Println("interval")
	})
	s.AddTimeoutFunc(time.Second*2, func() {
		log.Println("timeout")
	})
	s.Start()
	time.Sleep(time.Second * 4)
}

func TestFoo(t *testing.T) {
	c := make(chan struct{})

	go func() {
		close(c)
	}()

	<-c
	<-c
	<-c
	fmt.Println("after c")
}

func TestRunningDataRace(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			s.Start()
			wg.Done()
		}()
	}

	wg.Wait()
}
