package bot

import (
	"sync"
	"time"
)

// RateLimiter limits requests per user
type RateLimiter struct {
	mu           sync.Mutex
	requests     map[int64][]time.Time
	maxRequests  int
	windowDuration time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, windowDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:     make(map[int64][]time.Time),
		maxRequests:  maxRequests,
		windowDuration: windowDuration,
	}
}

// Allow checks if a request from userID is allowed
func (r *RateLimiter) Allow(userID int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.windowDuration)

	// Get existing requests for this user
	userRequests := r.requests[userID]

	// Filter out old requests
	var recentRequests []time.Time
	for _, reqTime := range userRequests {
		if reqTime.After(windowStart) {
			recentRequests = append(recentRequests, reqTime)
		}
	}

	// Check if under limit
	if len(recentRequests) >= r.maxRequests {
		return false
	}

	// Add new request
	recentRequests = append(recentRequests, now)
	r.requests[userID] = recentRequests

	return true
}

// Cleanup removes old entries from the rate limiter
func (r *RateLimiter) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.windowDuration)

	for userID, requests := range r.requests {
		var recentRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				recentRequests = append(recentRequests, reqTime)
			}
		}
		if len(recentRequests) == 0 {
			delete(r.requests, userID)
		} else {
			r.requests[userID] = recentRequests
		}
	}
}

