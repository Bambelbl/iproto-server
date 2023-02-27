package rate_limiter

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	logger   *log.Logger
	buckets  map[string]uint32
	scale    int64
	limit    uint32
	mutex    sync.RWMutex
	stopChan chan struct{}
}

const (
	DELETE_TIMEOUT = 1000
	INTERVAL_TIME  = 5000
)

func NewRateLimiter(loger *log.Logger, scale int64, limit uint32) *RateLimiter {
	rateLimiter := &RateLimiter{
		logger:   loger,
		buckets:  make(map[string]uint32),
		scale:    scale,
		limit:    limit,
		stopChan: make(chan struct{}),
	}
	go rateLimiter.removeOldLimiters()
	return rateLimiter
}

// Stop shutdown of RateLimiter
func (rl *RateLimiter) Stop() {
	rl.stopChan <- struct{}{}
}

// ValidRate checks if the client sends requests more than limit times per scale
func (rl *RateLimiter) ValidRate(IP string) bool {
	stamp := time.Now().UnixNano() / int64(time.Millisecond)
	bucketNumber := stamp / rl.scale
	key := IP + "_" + strconv.FormatInt(bucketNumber, 10)
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	if rps, exist := rl.buckets[key]; !exist {
		rl.buckets[key] = 1
		return true
	} else {
		if rps < rl.limit {
			rl.buckets[key]++
			return true
		} else {
			return false
		}
	}
}

// removeOldLimiters deletes old buckets times per interval
func (rl *RateLimiter) removeOldLimiters() {
	select {
	case <-rl.stopChan:
		return
	case <-time.Tick(time.Duration(INTERVAL_TIME)):
		stamp := time.Now().UnixNano() / int64(time.Millisecond)

		isOldTime := func(key string) bool {
			_, bucketTimeStr, _ := strings.Cut(key, "_")
			bucketTime, err := strconv.ParseInt(bucketTimeStr, 10, 64)
			if err != nil {
				rl.logger.Println("RateLimiter: parsing key error: ", err.Error())
			}
			return bucketTime < (stamp - DELETE_TIMEOUT)
		}

		rl.mutex.Lock()
		for key, _ := range rl.buckets {
			if isOldTime(key) {
				delete(rl.buckets, key)
			}
		}
		rl.mutex.Unlock()
	}
}
