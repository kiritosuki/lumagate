package main

import (
	"math/rand"
	"net/http"
	"time"
)

type Upstream struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

type PluginAuth struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
}

type PluginIPWhitelist struct {
	Enabled bool     `json:"enabled"`
	IPs     []string `json:"ips"`
}

type PluginRateLimit struct {
	Enabled   bool `json:"enabled"`
	WindowSec int  `json:"windowSec"`
	Max       int  `json:"max"`
}

type PluginCORS struct {
	Enabled  bool     `json:"enabled"`
	AllowAll bool     `json:"allowAll"`
	Origins  []string `json:"origins"`
}

type PluginTrafficSplit struct {
	Enabled   bool `json:"enabled"`
	V1Percent int  `json:"v1Percent"`
}

type Plugins struct {
	Auth         PluginAuth         `json:"auth"`
	IPWhitelist  PluginIPWhitelist  `json:"ipWhitelist"`
	RateLimit    PluginRateLimit    `json:"rateLimit"`
	CORS         PluginCORS         `json:"cors"`
	TrafficSplit PluginTrafficSplit `json:"trafficSplit"`
}

type Route struct {
	ID        string     `json:"id"`
	Prefix    string     `json:"prefix"`
	V1        []Upstream `json:"v1"`
	V2        []Upstream `json:"v2"`
	LBEnabled bool       `json:"lbEnabled"`
	Plugins   Plugins    `json:"plugins"`
}

var (
	routes  = map[string]*Route{}
	rrIndex = map[string]int{"v1": 0, "v2": 0}
)

func initState() {
	rand.Seed(time.Now().UnixNano())
	routes["users"] = &Route{
		ID:     "users",
		Prefix: "/api/users",
		V1: []Upstream{
			{Name: "A", URL: "http://serviceA_v1:9001", Weight: 1},
			{Name: "B", URL: "http://serviceB_v1:9002", Weight: 1},
			{Name: "C", URL: "http://serviceC_v1:9003", Weight: 1},
		},
		V2: []Upstream{
			{Name: "A", URL: "http://serviceA_v2:9011", Weight: 1},
			{Name: "B", URL: "http://serviceB_v2:9012", Weight: 1},
			{Name: "C", URL: "http://serviceC_v2:9013", Weight: 1},
		},
		LBEnabled: false,
		Plugins: Plugins{
			Auth:         PluginAuth{Enabled: false},
			IPWhitelist:  PluginIPWhitelist{Enabled: false},
			RateLimit:    PluginRateLimit{Enabled: false, WindowSec: 1, Max: 5},
			CORS:         PluginCORS{Enabled: false, AllowAll: true},
			TrafficSplit: PluginTrafficSplit{Enabled: false, V1Percent: 100},
		},
	}
}

func runPlugins(w http.ResponseWriter, r *http.Request, route *Route) bool {
	if !pluginAuth(w, r, route) {
		return false
	}
	if !pluginIPWhitelist(w, r, route) {
		return false
	}
	if !pluginRateLimit(w, r, route) {
		return false
	}
	if !pluginCORS(w, r, route) {
		return false
	}
	if !pluginTrafficSplit(route) {
		return false
	}
	return true
}
