package main

import (
	"math/rand"
)

func chooseUpstream(route *Route, group string) *Upstream {
	var ups []Upstream
	if group == "v1" {
		ups = route.V1
	} else {
		ups = route.V2
	}
	if len(ups) == 0 {
		return nil
	}
	if !route.LBEnabled {
		idx := rrIndex[group]
		u := ups[idx%len(ups)]
		rrIndex[group] = (idx + 1) % len(ups)
		return &Upstream{Name: u.Name, URL: u.URL, Weight: u.Weight}
	}
	total := 0
	for _, u := range ups {
		w := u.Weight
		if w <= 0 {
			w = 1
		}
		total += w
	}
	if total <= 0 {
		total = len(ups)
	}
	r := rand.Intn(total)
	acc := 0
	for _, u := range ups {
		w := u.Weight
		if w <= 0 {
			w = 1
		}
		acc += w
		if r < acc {
			return &Upstream{Name: u.Name, URL: u.URL, Weight: u.Weight}
		}
	}
	u := ups[0]
	return &Upstream{Name: u.Name, URL: u.URL, Weight: u.Weight}
}
