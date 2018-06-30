package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	dataFolder = "./data"

	versionStr = "v1.0"

	noteDayFormat  = "2006-01-02"
	noteTimeFormat = "15:04:05"
)

// NoteFile - File of Daily Notes
type NoteFile struct {
	Date  time.Time  `json:"Date"`
	Notes []NoteLine `json:"Notes"`
}

/**********************************************
			  NOTE LINE
*********************************************/

// Record - Convert to Record Line
func (nf *NoteFile) Record(line *NoteLine) []string {
	return []string{
		strconv.Itoa(line.Index),
		line.Note,
		line.Created.Format(noteTimeFormat),
		line.Modified.Format(noteTimeFormat),
	}
}

// FromRecord - Convert from Record
func (nf *NoteFile) FromRecord(record []string) (NoteLine, error) {
	var line NoteLine

	index64, err := strconv.ParseInt(record[0], 10, 32)
	if err != nil {
		return line, err
	}

	line.Index = int(index64)
	line.Note = record[1]

	line.Created, err = time.Parse(noteTimeFormat, record[2])
	if err != nil {
		return line, err
	}

	line.Modified, err = time.Parse(noteTimeFormat, record[3])
	if err != nil {
		return line, err
	}

	return line, nil
}

// Save - Save to File
func (nf *NoteFile) Save() error {
	dateStr := nf.Date.Format(noteDayFormat)

	f, err := os.Create(fmt.Sprintf("%s/%s.csv", dataFolder, dateStr))
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)

	// Header
	w.Write([]string{dateStr, versionStr})

	// Notes
	for _, nl := range nf.Notes {
		recLine := nf.Record(&nl)
		log.Println(recLine)

		err = w.Write(recLine)
		if err != nil {
			return err
		}
	}
	w.Flush()

	log.Printf("-- SAVED FILE -- ")

	return nil
}

// loadNoteFile - Load from File
func (nf *NoteFile) loadNoteFile(t time.Time) error {
	dateStr := t.Format(noteDayFormat)
	baseTime, err := time.Parse(noteDayFormat, dateStr)

	*nf = NoteFile{
		Date:  baseTime,
		Notes: []NoteLine{},
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s.csv", dataFolder, dateStr), os.O_RDONLY, 0666)
	if os.IsNotExist(err) {
		return os.MkdirAll(dataFolder, 0755)
	}
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Failed to open CSV \n %s", err.Error())
	}

	r := csv.NewReader(f)

	// Parse Header
	r.FieldsPerRecord = 2
	rec, err := r.Read()
	if err != nil {
		log.Printf("Failed to parse header \n %s", err.Error())
	} else {
		if rec[1] != versionStr {
			return fmt.Errorf("Version mismatch %s != %s", rec[1], versionStr)
		}

		ft, err := time.Parse(noteDayFormat, rec[0])
		if err != nil {
			return fmt.Errorf("Failed to parse Time \n %s", err.Error())
		}
		if ft != baseTime {
			return fmt.Errorf("Date mismatch %s != %s", ft.String(), baseTime.String())
		}
	}

	// Notes
	r.FieldsPerRecord = 4
	recList, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("Failed to readall \n %s", err.Error())
	}

	nf.Notes = make([]NoteLine, len(recList), len(recList))

	for i, nl := range recList {
		nf.Notes[i], err = nf.FromRecord(nl)
		if err != nil {
			return fmt.Errorf("Failed to do From Record %d \n %s \n %s", i, nl, err.Error())
		}
	}

	return nil
}

// New - Create New Note
func (nf *NoteFile) New(note string, t time.Time) (NoteLine, error) {
	// Sanity Checks
	if note == "" {
		return NoteLine{}, fmt.Errorf("Empty Note")
	}

	// Check Date
	if (t.Before(nf.Date)) || t.After(nf.Date.Add(time.Hour*24)) {
		err := nf.loadNoteFile(t)
		if err != nil {
			return NoteLine{}, err
		}
	}

	line := NoteLine{
		len(nf.Notes),
		note,
		t,
		t,
	}

	nf.Notes = append(nf.Notes, line)

	// Save Change
	err := nf.Save()
	if err != nil {
		return line, err
	}

	return line, nil
}

// Edit - Edit exsisting Note
func (nf *NoteFile) Edit(i int, note string, t time.Time) (NoteLine, error) {
	// Check Date
	if (t.Before(nf.Date)) || t.After(nf.Date.Add(time.Hour*24)) {
		err := nf.loadNoteFile(t)
		if err != nil {
			return NoteLine{}, err
		}
	}

	// Sanity Checks
	if i >= len(nf.Notes) {
		return NoteLine{}, fmt.Errorf("%d out of range", i)
	}

	line := &(nf.Notes[i])

	line.Note = note
	line.Modified = t

	// Save Change
	err := nf.Save()
	if err != nil {
		return *line, err
	}

	return *line, nil
}

// Get - Get Note from Day and Index
func (nf *NoteFile) Get(t time.Time, i int) (NoteLine, error) {
	// Check Date
	if (t.Before(nf.Date)) || t.After(nf.Date.Add(time.Hour*24)) {
		err := nf.loadNoteFile(t)
		if err != nil {
			return NoteLine{}, err
		}
	}

	// Sanity Checks
	if i >= len(nf.Notes) {
		return NoteLine{}, fmt.Errorf("%d out of range", i)
	}

	return nf.Notes[i], nil
}

// GetDay - Get Full Days Notes
func (nf *NoteFile) GetDay(t time.Time) []NoteLine {
	// Check Date
	if (t.Before(nf.Date)) || t.After(nf.Date.Add(time.Hour*24)) {
		err := nf.loadNoteFile(t)
		if err != nil {
			return []NoteLine{}
		}
	}

	return nf.Notes
}
