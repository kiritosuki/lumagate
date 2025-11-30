package main

import (
	"net/http"
	"sync"
	"time"
)

type limiterState struct {
	windowStart time.Time
	count       int
}

var (
	limitMu     sync.Mutex
	limitsState = map[string]*limiterState{}
)

func pluginRateLimit(w http.ResponseWriter, r *http.Request, route *Route) bool {
	c := route.Plugins.RateLimit
	if !c.Enabled {
		return true
	}
	limitMu.Lock()
	st := limitsState[route.ID]
	now := time.Now()
	if st == nil {
		st = &limiterState{windowStart: now, count: 0}
		limitsState[route.ID] = st
	}
	if now.Sub(st.windowStart) > time.Duration(c.WindowSec)*time.Second {
		st.windowStart = now
		st.count = 0
	}
	if st.count >= c.Max {
		limitMu.Unlock()
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return false
	}
	st.count++
	limitMu.Unlock()
	return true
}
