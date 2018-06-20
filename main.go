package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

var (
	logBuffer bytes.Buffer

	serverSettings string
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

	// Setup Console
	cmdChan := make(chan string)
	go consoleReadLoop(cmdChan)

	// Setup Website
	servRouter := launchWeb(serverSettings)

	servRouter.HandleFunc("/dashboard", handleAdmin)
	servRouter.HandleFunc("/new", handleNoteNew)
	servRouter.HandleFunc("/edit", handleNoteEdit)
	servRouter.HandleFunc("/today", handleNoteToday)
	servRouter.HandleFunc("/past", handleNotePast)

	isRunning := true
	for isRunning {
		select {
		case cmdLine := <-cmdChan:
			if cmdLine == "quit" {
				isRunning = false
				break
			}
		}
	}

	log.Printf("Shutting down")
}

func consoleReadLoop(c chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd == "panic" {
			panic("Forced Panic")
		}

		select {
		case c <- cmd:
		default: // non blocking
		}
	}
}
