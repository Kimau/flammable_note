package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func handleNoteNew(w http.ResponseWriter, req *http.Request) {

	corsMe(w, req)

}

func handleNoteEdit(w http.ResponseWriter, req *http.Request) {
	corsMe(w, req)

}

func handleNoteToday(w http.ResponseWriter, req *http.Request) {
	corsMe(w, req)

}

func handleNotePast(w http.ResponseWriter, req *http.Request) {
	corsMe(w, req)

}
