package main

import "net/http"

func pluginAuth(w http.ResponseWriter, r *http.Request, route *Route) bool {
	c := route.Plugins.Auth
	if !c.Enabled {
		return true
	}
	if r.Header.Get("X-API-Key") != c.Key {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}
