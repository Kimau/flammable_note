package main

import (
	"time"
)

const ()

// NoteLine - Single Line in Note
type NoteLine struct {
	Index    int       `json:"ID"`
	Note     string    `json:"Note"`
	Created  time.Time `json:"Created"`
	Modified time.Time `json:"Modified"`
}

// NoteStore - Place to store notes
type NoteStore interface {
	Save() error

	New(note string, t time.Time) (NoteLine, error)
	Edit(i int, note string, t time.Time) (NoteLine, error)

	Get(t time.Time, i int) (NoteLine, error)
	GetDay(t time.Time) []NoteLine
}
