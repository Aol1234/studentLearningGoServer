package sessions

import (
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

var c = cache.New(240*time.Minute, 40*time.Minute) // Start cache

func Get(token string) (interface{}, bool) {
	// Get UserId using session token
	value, found := c.Get(token)
	if found == false {
		log.Println("Invalid token")
		return nil, false
	}
	return value, true
}

func Set(token string, UserId uint, expiry time.Duration) {
	// Set session
	c.Set(token, UserId, expiry)
}

func SetSession(UserId uint, bearerToken string) bool {
	// Set session
	if bearerToken == "" {
		return false
	}
	Set("Bearer  "+bearerToken, UserId, 0)
	return true
}
