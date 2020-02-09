package telegrabotlib

import (
	"gopkg.in/tucnak/telebot.v2"
	"testing"
	"time"
)

func TestScheduler_run(t *testing.T) {
	s := NewScheduler()

	modifyA := 1
	modifyB := 2

	s.AddToSchedule(ScheduleTask{
		Interval: time.Second,
		Fn: func(bot *telebot.Bot) {
			modifyA = 2
		},
	})

	s.AddToSchedule(ScheduleTask{
		Interval: time.Second,
		Fn: func(bot *telebot.Bot) {
			modifyB = 1
		},
	})

	go func() {
		t := time.NewTicker(3 * time.Second)

		select {
		case <-t.C:
			s.Stop()
		}
	}()

	s.run(nil)

	if modifyA != 2 {
		t.Error("Task 1 wasn't executed")
	}

	if modifyB != 1{
		t.Error("Task 2 wasn't executed")
	}
}