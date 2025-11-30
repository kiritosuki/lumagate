package main

import (
	"net"
	"net/http"
)

func pluginIPWhitelist(w http.ResponseWriter, r *http.Request, route *Route) bool {
	c := route.Plugins.IPWhitelist
	if !c.Enabled {
		return true
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	allowed := false
	for _, a := range c.IPs {
		if a == ip {
			allowed = true
			break
		}
	}
	if !allowed {
		http.Error(w, "forbidden", http.StatusForbidden)
		return false
	}
	return true
}
