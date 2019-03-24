package sessions

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var c = cache.New(240*time.Minute, 40*time.Minute)

func Get(token string) (interface{}, bool) {
	value, found := c.Get(token)
	if found == false {
		panic("Invalid token")
		return nil, false
	}
	return value, true
}

func Set(token string, UserId uint, expiry time.Duration) {
	c.Set(token, UserId, expiry)
}
