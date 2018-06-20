package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "net/http/pprof"
)

// GZipJSON - Zips and Sends JSON
func GZipJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		err := json.NewEncoder(gz).Encode(data)
		if err != nil {
			errMsg := fmt.Sprintf("Encoding Error with GZIP: %s", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, 511)
		}
		gz.Close()
	} else {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			errMsg := fmt.Sprintf("Encoding Error %s", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, 511)
		}
	}
}
