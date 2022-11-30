package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// read from file
			data, err := ioutil.ReadFile("data.json")
			if err != nil {
				// return error in response
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
				return
			}

			// set response
			w.Write(data)
		case "POST":
			// read POST request body
			data, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "Failed to get data: %s", err.Error())
				return
			}

			// read existing data from file
			existing_data, err := ioutil.ReadFile("data.json")
			if err != nil {
				fmt.Fprintf(w, "Failed to get existing data: %s", err.Error())
			}

			// parse existing string data to json
			var todos []string
			if err = json.Unmarshal(existing_data, &todos); err != nil {
				fmt.Fprintf(w, "Failed to parse existing data: %s", err.Error())
				return
			}

			// add new data
			todos = append(todos, string(data))
			// convert parsed data to string
			output, err := json.Marshal(todos)
			if err != nil {
				fmt.Fprintf(w, "Failed to convert to json: %s", err.Error())
			}

			// write to file
			if err = ioutil.WriteFile("data.json", output, 0666); err != nil {
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
