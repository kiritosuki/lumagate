package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func adminRoutesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := make([]*Route, 0, len(routes))
		for _, rt := range routes {
			list = append(list, rt)
		}
		writeJSON(w, list)
	case http.MethodPost:
		var rt Route
		if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if rt.ID == "" {
			rt.ID = rt.Prefix
		}
		routes[rt.ID] = &rt
		writeJSON(w, map[string]string{"status": "ok"})
	default:
		http.Error(w, "method", http.StatusMethodNotAllowed)
	}
}

func adminRoutesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var rt Route
	if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if rt.ID == "" {
		http.Error(w, "id", http.StatusBadRequest)
		return
	}
	routes[rt.ID] = &rt
	writeJSON(w, map[string]string{"status": "ok"})
}

func adminRoutesIDHandler(w http.ResponseWriter, r *http.Request) {
	// /admin/routes/:id
	parts := splitPath(r.URL.Path)
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	if r.Method == http.MethodDelete {
		delete(routes, id)
		writeJSON(w, map[string]string{"status": "ok"})
		return
	}
	http.Error(w, "method", http.StatusMethodNotAllowed)
}

func adminLogsHandler(w http.ResponseWriter, r *http.Request) {
	tail := 100
	if v := r.URL.Query().Get("tail"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			tail = n
		}
	}
	p := filepath.Join("logs", "access.log")
	f, err := os.Open(p)
	if err != nil {
		http.Error(w, "no log", http.StatusNotFound)
		return
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	lines := splitLines(string(b))
	if len(lines) > tail {
		lines = lines[len(lines)-tail:]
	}
	writeJSON(w, lines)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func splitPath(p string) []string {
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")
	var out []string
	for _, s := range parts {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, "\n")
}
