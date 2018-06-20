package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func getTimeFormRequest(vars map[string]string) (time.Time, error) {
	// Setup Timestamp
	noteTime := time.Now()
	noteTimeStr, ok := vars["timestamp"]
	if ok {
		timeSec, err := strconv.ParseInt(noteTimeStr, 10, 64)
		if err != nil {
			return noteTime, err
		}
		noteTime = time.Unix(timeSec, 0)
	}

	return noteTime, nil
}

func handleAdmin(w http.ResponseWriter, req *http.Request) {
	corsMe(w, req)

	// WARNING :: DUMMY Auth
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Failed Auth: ", err)
		return
	}

	fmt.Fprintln(w, "Cool")
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Root")
}

func handleNoteNew(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	if req.Method != "POST" {
		http.Error(w, "New Notes must be by POST", 400)
		return
	}

	corsMe(w, req)

	// Time
	noteTime, err := getTimeFormRequest(vars)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Read Body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// New Note
	line, err := noteStore.New(string(body), noteTime)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Response
	GZipJSON(w, req, line)
}

func handleNoteEdit(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	if req.Method != "POST" {
		http.Error(w, "Edit Notes must be by POST", 400)
	}

	corsMe(w, req)

	// Index
	indexStr, ok := vars["index"]
	if !ok {
		http.Error(w, "Must give index", 400)
		return
	}
	i, err := strconv.ParseInt(indexStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Time
	noteTime, err := getTimeFormRequest(vars)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Read Body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Edit Note
	line, err := noteStore.Edit(int(i), string(body), noteTime)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Response
	GZipJSON(w, req, line)
}

func handleNoteToday(w http.ResponseWriter, req *http.Request) {
	corsMe(w, req)

	GZipJSON(w, req, noteStore.GetDay(time.Now()))
}

func handleNotePast(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	corsMe(w, req)

	// Time
	noteTime, err := getTimeFormRequest(vars)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	GZipJSON(w, req, noteStore.GetDay(noteTime))
}
