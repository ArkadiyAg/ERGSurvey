package main

import (
	"ERGSurvey/back/app/survey"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := "8080"

	surv := survey.CreateDummySurvey()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderQuestion(w, "templates/survey.component.gohtml", surv.CurrentQuestion())
	})

	http.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Questions = %v\n", surv)
		renderAnswers(w, "templates/table.component.gohtml", surv)
	})

	http.HandleFunc("/question", func(w http.ResponseWriter, r *http.Request) {
		renderNewQuestion(w, "templates/question.component.gohtml")
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Printf("Method submit has been called\n")
			res := &survey.ResponsePayload{}
			err := readJSON(w, r, res)
			if err != nil {
				fmt.Printf("Falsed to parse JSON: %v\n", err)
			}
			ip := readUserIP(r)
			fmt.Printf("Request from: %s\n", ip)
			surv.Increment(res.Id, ip)
		}
	})

	http.HandleFunc("/newQuestion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Printf("Method submit has been called\n")
			res := &survey.NewQuestionPayload{}
			err := readJSON(w, r, res)
			if err != nil {
				fmt.Printf("Falsed to parse JSON: %v\n", err)
			}

			if res.Pin == "31415926" {
				surv.AddQuestion(res.Name, res.Q1, res.Q2, res.Q3)
			}
		}
	})

	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Printf("Method latest has been called\n")
			data, err := json.Marshal(*surv.CurrentQuestion())
			if err != nil {
				fmt.Printf("failed to encode the object to JSON: %v", err)
			}
			w.Write(data)
		}
	})

	http.HandleFunc("/setQuestion", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Get the query parameters from the request
			params := r.URL.Query()

			// Get the "name" parameter value and print it
			pin := params.Get("pin")
			num := params.Get("num")
			//fmt.Printf("Name parameter value is: %s\n", name)

			if pin == "31415926" {
				numInt, _ := strconv.Atoi(num)
				surv.SetQuestion(numInt)
			}
		}
	})

	fmt.Printf("Starting survey frontend on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panic(err)
	}
}
