package scheduler

import (
	"log"

	"gorm.io/gorm"
)

func CleanRequestLogs(DB *gorm.DB) {
	if err := DB.Exec("DELETE FROM request_logs WHERE request_time < NOW() - INTERVAL '30 days'").Error; err != nil {
		log.Println("Error cleaning request logs: ", err)
	}
}
