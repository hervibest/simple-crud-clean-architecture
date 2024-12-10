package helper

import (
	"reflect"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

var (
	RateLimiter = make(map[string]int64)
	PointLimit  = make(map[string]int64)
	Mutex       sync.Mutex
	MutexMap    = make(map[string]*sync.Mutex)
)

func SanitiseStruct(input interface{}) {
	val := reflect.ValueOf(input).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			// Sanitasi field string
			sanitised := bluemonday.UGCPolicy().Sanitize(field.String())
			field.SetString(sanitised)
		}
	}
}

func RateLimit(key string, maxAttempt, maxTimeInSeconds int) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	now := time.Now().Unix()

	if v, ok := RateLimiter[key]; ok && now-v < int64(maxTimeInSeconds) {
		RateLimiter[key]++
		PointLimit[key]++
		if PointLimit[key] > int64(maxAttempt) {
			return false
		}
	} else {
		RateLimiter[key] = now
		PointLimit[key] = 1
	}

	return true
}

func LockThread(key string) {
	mutex, ok := MutexMap[key]
	if !ok {
		mutex = &sync.Mutex{}
		MutexMap[key] = mutex
	}
	mutex.Lock()
}

func UnlockThread(key string) {
	mutex, ok := MutexMap[key]
	if ok {
		mutex.Unlock()
	}
}
