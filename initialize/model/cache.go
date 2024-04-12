package model

import (
	"sync"
	"time"
)

// UserCache /**
type UserCache struct {
	JWT         string
	OrgID       string
	LastUpdated time.Time
}

var cache = struct {
	sync.RWMutex
	items map[string]UserCache
}{items: make(map[string]UserCache)}

func SetJWT(rt string, jwt string) {
	cache.Lock()
	userCache := cache.items[rt]
	userCache.JWT = jwt
	userCache.LastUpdated = time.Now()
	cache.items[rt] = userCache
	cache.Unlock()
}

func SetOrgID(rt string, orgID string) {
	cache.Lock()
	userCache := cache.items[rt]
	userCache.OrgID = orgID
	userCache.LastUpdated = time.Now()
	cache.items[rt] = userCache
	cache.Unlock()
}

func GetUserCache(rt string) (UserCache, bool) {
	cache.RLock()
	userCache, exists := cache.items[rt]
	cache.RUnlock()
	if !exists || time.Since(userCache.LastUpdated) > time.Hour { // Assume cache valid for 1 hour
		return UserCache{}, false
	}
	return userCache, true
}
