package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			log.Info("Received get requests")

			// read from file
			data, err := ioutil.ReadFile("data.json")
			if err != nil {
				// return error in response
				log.WithError(err).Error("Failed to read from data.json")
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
				return
			}

			// set response
			w.Write(data)
		case "POST":
			log.Info("Received post requests")

			// read POST request body
			data, err := io.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Error("Failed to read data from post")
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
				return
			}

			// read existing data from file
			existing_data, err := ioutil.ReadFile("data.json")
			if err != nil {
				log.WithError(err).Error("Failed to read from data.json")
				fmt.Fprintf(w, "Failed to get existing data: %s", err.Error())
			}

			// parse existing string data to json
			var todos []string
			if err = json.Unmarshal(existing_data, &todos); err != nil {
				log.WithError(err).Error("Failed to parse existing data")
				fmt.Fprintf(w, "Failed to parse existing data: %s", err.Error())
				return
			}

			// add new data
			todos = append(todos, string(data))
			// convert parsed data to string
			output, err := json.Marshal(todos)
			if err != nil {
				log.WithError(err).WithField("todos", todos).Error("Failed to convert todos to json")
				fmt.Fprintf(w, "Failed to convert to json: %s", err.Error())
			}

			// write to file
			if err = ioutil.WriteFile("data.json", output, 0666); err != nil {
				log.WithError(err).WithField("data", output).Error("Failed to write data to data.json")
				fmt.Fprintf(w, "Failed to write to data.json: %s", err.Error())
			}

			// set response
			w.Write([]byte("Successfully wrote data"))
		default:
			fmt.Fprintf(w, "Method %s not supported", r.Method)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server crashed: %s", err.Error())
	}
}
