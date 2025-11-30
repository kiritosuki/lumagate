package main

import (
    "log"
    "net/http"
)

func main() {
    initState()
    http.HandleFunc("/admin/routes", adminRoutesHandler)
    http.HandleFunc("/admin/routes/update", adminRoutesUpdateHandler)
    http.HandleFunc("/admin/routes/", adminRoutesIDHandler)
    http.HandleFunc("/admin/logs", adminLogsHandler)
    http.HandleFunc("/", gatewayHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

