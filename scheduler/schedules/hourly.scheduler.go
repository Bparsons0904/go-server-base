package scheduler

import (
	"github.com/bparsons094/go-server-base/utils"
)

func CleanUserCache() {
	utils.ClearExpiredUsers()
}
