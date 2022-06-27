package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mehmet.com/cors"
	"net/http"
	"strings"
)

const personsPath = "persons"

func PersonsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		personList, err := getPersonList()

		personListJSON, err := json.Marshal(personList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(personListJSON)
		if err != nil {
			log.Fatal("Error!")
		}
	case http.MethodPost:
		newPerson := Person{}
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &newPerson)
		w.Header().Set("Content-Type", "application/json")

		err = addPerson(&newPerson)
		if err != nil {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(201)
		}

	}
}

func PersonHandler(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/persons/")
	switch r.Method {
	case http.MethodGet:

		person, err := getPerson(id)

		personJSON, err := json.Marshal(person)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(personJSON)
		if err != nil {
			log.Fatal("Error!")
		}
	case http.MethodPut:
		newPerson := Person{}
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &newPerson)
		w.Header().Set("Content-Type", "application/json")

		err = updatePerson(&newPerson, id)
		if err != nil {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}
	case http.MethodDelete:

		count, err := deletePerson(id)

		if count < 1 {
			w.WriteHeader(404)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(nil)
		if err != nil {
			log.Fatal("Error!")
		}

	}
}

func SetupRoutes(apiBasePath string) {
	productsHandler := http.HandlerFunc(PersonsHandler)
	productHandler := http.HandlerFunc(PersonHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, personsPath), cors.Middleware(productsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, personsPath), cors.Middleware(productHandler))
}
