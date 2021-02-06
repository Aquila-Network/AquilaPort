package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var myRouter = mux.NewRouter().StrictSlash(true)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func existsDatabaseRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	if existsDatabase(databaseName) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func getDatabaseRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	if existsDatabase(databaseName) {
		json.NewEncoder(w).Encode(map[string]string{
			"db_name": databaseName,
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func createNewDatabaseRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	status := createNewDatabase(databaseName)

	if status {
		json.NewEncoder(w).Encode(map[string]bool{
			"ok": true,
		})
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func createNewDocumentRouter(w http.ResponseWriter, r *http.Request) {
	// decode json body
	var documents []Document

	json.NewEncoder(w).Encode(documents)
}

func deleteDocumentRouter(w http.ResponseWriter, r *http.Request) {
	var ids []string

	json.NewEncoder(w).Encode(ids)
}

func handleRequests(port string) {
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/{databaseName}", existsDatabaseRouter).Methods("HEAD")
	myRouter.HandleFunc("/{databaseName}", getDatabaseRouter).Methods("GET")
	myRouter.HandleFunc("/{databaseName}", createNewDatabaseRouter).Methods("PUT")
	myRouter.HandleFunc("/create", createNewDocumentRouter).Methods("POST")
	myRouter.HandleFunc("/delete", deleteDocumentRouter).Methods("POST")

	// Run server
	fmt.Println("Aquila Port running at localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}
