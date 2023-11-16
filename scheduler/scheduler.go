package scheduler

import (
	"log"
	"time"

	scheduler "github.com/bparsons094/go-server-base/scheduler/schedules"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func InitScheduler(DB *gorm.DB) {
	log.Println("Starting Schedules")

	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Hour().Do(func() {
		scheduler.CleanUserCache()
	})

	s.Every(1).Day().At("00:00").Do(func() {
		scheduler.CleanRequestLogs(DB)
	})
	s.StartBlocking()
}
