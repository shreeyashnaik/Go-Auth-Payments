package cache

import (
	"context"
	"time"

	"github.com/Shreeyash-Naik/Go-Auth/common/db"
)

var (
	CacheClient = db.RedisDB
	Ctx         = context.Background()
)

func SetValueEx(key string, value int, exp time.Duration) {
	CacheClient.SetEX(Ctx, key, value, exp)
}
