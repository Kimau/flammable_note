package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

var (
	logBuffer bytes.Buffer

	serverSettings string
	noteStore      NoteFile
)

func init() {
	serverSettings = "server.json"
}

func main() {
	// Setup Log
	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(f, &logBuffer, os.Stdout))

	// Load Data
	err = noteStore.loadNoteFile(time.Now())
	if err != nil {
		panic(err)
	}

	// Setup Website
	servRouter := launchWeb(serverSettings)

	servRouter.HandleFunc("/", handleRoot)
	servRouter.HandleFunc("/dashboard", handleAdmin)
	servRouter.HandleFunc("/new", handleNoteNew).Methods("POST")
	servRouter.HandleFunc("/edit/{index}", handleNoteEdit).Methods("POST")
	servRouter.HandleFunc("/today", handleNoteToday)
	servRouter.HandleFunc("/past", handleNotePast)

	ticker := time.NewTicker(time.Minute)

	isRunning := true
	for isRunning {
		select {
		case <-ticker.C:
			// Check It
		}
	}

	log.Printf("Shutting down")
}
