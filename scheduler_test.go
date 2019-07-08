package scheduler

import (
	"log"
	"testing"
	"time"
)

var s = New()

func TestScheduler_AddIntervalFunc(t *testing.T) {

	s.AddIntervalFunc(time.Second, func() {
		log.Println("interval")
	})
	s.AddTimeoutFunc(time.Second * 2, func() {
		log.Println("timeout")
	})
	s.AddTimeoutFunc(time.Second * 2, func() {
		log.Println("timeout")
	})
	s.Stop()
	time.Sleep(time.Second * 4)
	s.Stop()
	time.Sleep(time.Second * 2)

}

