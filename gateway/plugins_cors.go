package main

import "net/http"

func pluginCORS(w http.ResponseWriter, r *http.Request, route *Route) bool {
	c := route.Plugins.CORS
	if !c.Enabled {
		return true
	}
	if c.AllowAll {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return false
	}
	return true
}
