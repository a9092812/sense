package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/Kartik30R/sense/config"
	"golang.org/x/time/rate"
)

type deviceLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	globalLimiter *rate.Limiter

	deviceLimiters map[string]*deviceLimiter
	deviceRPS      rate.Limit
	burst          int

	mutex sync.RWMutex
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {

	rl := &RateLimiter{
		globalLimiter: rate.NewLimiter(
			rate.Limit(cfg.GlobalRPS),
			cfg.BurstSize,
		),

		deviceLimiters: make(map[string]*deviceLimiter),
		deviceRPS:      rate.Limit(cfg.DeviceRPS),
		burst:          cfg.BurstSize,
	}

	go rl.cleanupLoop()

	return rl
}

func (r *RateLimiter) getDeviceLimiter(deviceID string) *rate.Limiter {

	// fast path
	r.mutex.RLock()
	dl, exists := r.deviceLimiters[deviceID]
	r.mutex.RUnlock()

	if exists {
		dl.lastSeen = time.Now()
		return dl.limiter
	}

	// slow path
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// double check
	dl, exists = r.deviceLimiters[deviceID]
	if exists {
		dl.lastSeen = time.Now()
		return dl.limiter
	}

	limiter := rate.NewLimiter(r.deviceRPS, r.burst)

	r.deviceLimiters[deviceID] = &deviceLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

func (r *RateLimiter) Allow(deviceID string) bool {

	if !r.globalLimiter.Allow() {
		return false
	}

	deviceLimiter := r.getDeviceLimiter(deviceID)

	return deviceLimiter.Allow()
}

func (r *RateLimiter) Wait(ctx context.Context, deviceID string) error {

	if err := r.globalLimiter.Wait(ctx); err != nil {
		return err
	}

	deviceLimiter := r.getDeviceLimiter(deviceID)

	return deviceLimiter.Wait(ctx)
}

func (r *RateLimiter) cleanupLoop() {

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.cleanup()
	}
}

func (r *RateLimiter) cleanup() {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	expiration := 30 * time.Minute

	for deviceID, dl := range r.deviceLimiters {
		if time.Since(dl.lastSeen) > expiration {
			delete(r.deviceLimiters, deviceID)
		}
	}
}
