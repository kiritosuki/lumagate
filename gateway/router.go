package main

import (
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func gatewayHandler(w http.ResponseWriter, r *http.Request) {
	route := matchRoute(r.URL.Path)
	if route == nil {
		http.NotFound(w, r)
		return
	}
	if !runPlugins(w, r, route) {
		return
	}
	group := chooseGroup(route)
	upstream := chooseUpstream(route, group)
	if upstream == nil {
		http.Error(w, "no upstream", http.StatusBadGateway)
		return
	}
	target, _ := url.Parse(upstream.URL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		http.Error(rw, "upstream error", http.StatusBadGateway)
	}
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		p := strings.TrimPrefix(r.URL.Path, route.Prefix)
		if p == "" {
			p = "/"
		}
		req.URL.Path = p
	}
	start := time.Now()
	sr := &statusRecorder{ResponseWriter: w}
	proxy.ServeHTTP(sr, r)
	duration := time.Since(start)
	clientIP := clientIPFromReq(r)
	writeAccessLog(r.Method, r.URL.Path, sr.status, upstream.URL, duration, clientIP)
}

func matchRoute(path string) *Route {
	var best *Route
	for _, rt := range routes {
		if strings.HasPrefix(path, rt.Prefix) {
			if best == nil || len(rt.Prefix) > len(best.Prefix) {
				best = rt
			}
		}
	}
	return best
}

func chooseGroup(route *Route) string {
	if route.Plugins.TrafficSplit.Enabled {
		p := route.Plugins.TrafficSplit.V1Percent
		if p <= 0 {
			return "v2"
		}
		if p >= 100 {
			return "v1"
		}
		if rand.Intn(100) < p {
			return "v1"
		}
		return "v2"
	}
	return "v1"
}

func clientIPFromReq(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
