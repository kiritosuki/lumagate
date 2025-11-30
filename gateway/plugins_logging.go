package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var logMu sync.Mutex

func writeAccessLog(method, path string, status int, upstream string, d time.Duration, ip string) {
	p := filepath.Join("logs", "access.log")
	os.MkdirAll(filepath.Dir(p), 0755)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	logMu.Lock()
	fmt.Fprintf(f, "%s %s %s %d %s %dms %s\n", time.Now().Format(time.RFC3339), ip, method, status, upstream, d.Milliseconds(), path)
	logMu.Unlock()
}
