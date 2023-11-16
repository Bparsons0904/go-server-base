package scheduler

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func InitScheduler(DB *gorm.DB) {
	log.Println("Starting Schedules")

	s := gocron.NewScheduler(time.UTC)

	s.StartBlocking()
}
