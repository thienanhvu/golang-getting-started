package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const fileName = "data.json"

type Note struct {
	Id   int    `json:"id"`
	Note string `json:"note"`
}

func checkNoteExists(id int) (bool, error) {
	notes, err := getNotesFromFile()
	for _, note := range notes {
		if note.Id == id {
			return true, err
		}
	}
	return false, err
}

func getNotesFromFile() ([]Note, error) {
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var notes []Note
	err = json.Unmarshal(fileData, &notes)

	return notes, nil
}

func updateAllNotes(notes []Note) error {
	output, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("data.json", output, 0666)
	return err
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	log.Infof("GET %s\n", r.URL)

	notes, err := getNotesFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, _ := json.Marshal(notes)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func addNewNotes(w http.ResponseWriter, r *http.Request) {
	log.Infof("POST %s\n", r.URL)

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isExist, err := checkNoteExists(note.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if isExist {
		http.Error(w, "id already taken", http.StatusBadRequest)
		return
	}

	notes, err := getNotesFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notes = append(notes, note)
	err = updateAllNotes(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully wrote data"))
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	log.Infof("PUT %s\n", r.URL)

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isExist, err := checkNoteExists(note.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isExist {
		http.Error(w, "id not exists", http.StatusBadRequest)
		return
	}
	notes, err := getNotesFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i := range notes {
		if notes[i].Id == note.Id {
			notes[i].Note = note.Note
			break
		}
	}
	err = updateAllNotes(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Successfully update data"))
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	log.Infof("DELETE %s\n", r.URL)
	
	params := r.URL.Query()
	fmt.Printf("%v", params)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllNotes(w, r)
		case "POST":
			addNewNotes(w, r)
		case "PUT":
			updateNote(w, r)
		case "DELETE":
			deleteNote(w, r)
		default:
			fmt.Printf("Method %s not supported\n", r.Method)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server crashed: %s", err.Error())
	}
}
