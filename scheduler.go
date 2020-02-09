package telegrabotlib

import (
    "gopkg.in/tucnak/telebot.v2"
    "log"
    "time"
)

type Scheduler struct {
    list []ScheduleTask
    stop chan struct{}
}

func NewScheduler() *Scheduler {
    return &Scheduler{
        stop: make(chan struct{}),
    }
}

func (d *Scheduler) AddToSchedule(task ScheduleTask) {
    d.list = append(d.list, task) 
}

func (d *Scheduler) Stop() {
    close(d.stop)
}

type ScheduleTask struct {
    Fn ScheduleFn
    Interval time.Duration 
}
type ScheduleFn func(bot *telebot.Bot)

func (d *Scheduler) run(bot *telebot.Bot) {
    if len(d.list) == 0 {
        log.Print("No one job scheduled")
        return
    }

    for _, r := range d.list {
        go func(t ScheduleTask, bot *telebot.Bot) {
            ticker := time.NewTicker(t.Interval)

            for {
                select {
                case <-ticker.C:
                    t.Fn(bot)
                }
            }
        }(r, bot)
    }

    <-d.stop
}
