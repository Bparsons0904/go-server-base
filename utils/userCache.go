package utils

import (
	"sync"
	"time"

	"github.com/bparsons094/go-server-base/models"
	"github.com/google/uuid"
)

var userCache = UserCache{
	users: make(map[uuid.UUID]CacheItem),
}

type CacheItem struct {
	User      models.User
	ExpiresAt time.Time
}

type UserCache struct {
	users map[uuid.UUID]CacheItem
	mutex sync.RWMutex
}

func GetUser(id uuid.UUID) (models.User, bool) {
	userCache.mutex.RLock()
	defer userCache.mutex.RUnlock()

	cacheItem, found := userCache.users[id]
	if !found {
		return models.User{}, false
	}

	if time.Now().After(cacheItem.ExpiresAt) {
		return models.User{}, false
	}

	return cacheItem.User, true
}

func SetUser(id uuid.UUID, user models.User) {
	userCache.mutex.Lock()
	defer userCache.mutex.Unlock()

	expiration := time.Now().Add(30 * time.Minute)
	userCache.users[id] = CacheItem{
		User:      user,
		ExpiresAt: expiration,
	}
}

func ClearExpiredUsers() {
	userCache.mutex.Lock()
	defer userCache.mutex.Unlock()
	for id, cacheItem := range userCache.users {
		if time.Now().After(cacheItem.ExpiresAt) {
			delete(userCache.users, id)
		}
	}
}
