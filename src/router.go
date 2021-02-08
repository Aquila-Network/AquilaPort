package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var myRouter = mux.NewRouter().StrictSlash(true)

func nodeInfoRouter(w http.ResponseWriter, r *http.Request) {
	uuid := getPortUUID()
	json.NewEncoder(w).Encode(map[string]string{
		"nodeId": uuid,
	})
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

func getAllDocumentsRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	documents := getAllDocuments(databaseName)

	if documents != nil {
		json.NewEncoder(w).Encode(documents)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func createBulkDocumentsRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	// decode json body
	var documents []Document
	json.NewDecoder(r.Body).Decode(&documents)

	status := createNewDocuments(databaseName, documents)

	if status != nil {
		json.NewEncoder(w).Encode(status)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getDocumentChangesRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	changes := getDocumentChanges(databaseName, "0")

	json.NewEncoder(w).Encode(changes)
}

func getRevDiffRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	// decode change map
	changeMap := make(map[string][]string)
	json.NewDecoder(r.Body).Decode(&changeMap)

	diff := getRevDiff(databaseName, changeMap)

	json.NewEncoder(w).Encode(diff)
}

func ensureCommitRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	status := existsDatabase(databaseName)

	if status {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]bool{
			"ok": true,
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func replCPointRecordRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	// decode json body
	var rcpoint ReplCheckpoint
	json.NewDecoder(r.Body).Decode(&rcpoint)

	status := replCPointRecord(databaseName, rcpoint)

	if status {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rcpoint)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getReplCPointRecordRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	rcpointID := params["replID"]

	// decode json body
	var rcpoint ReplCheckpoint
	json.NewDecoder(r.Body).Decode(&rcpoint)

	status, replRecord := getReplCPointRecord(databaseName, rcpointID)

	if status {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(replRecord)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func enableReplicationRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	// decode json body
	var replReq ReplRequest
	json.NewDecoder(r.Body).Decode(&replReq)

	status := updateReplicationStatus(databaseName, replReq, "active")

	if status {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(replReq)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func disableReplicationRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	// decode json body
	var replReq ReplRequest
	json.NewDecoder(r.Body).Decode(&replReq)

	status := updateReplicationStatus(databaseName, replReq, "disabled")

	if status {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(replReq)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleRequests(port string) {
	myRouter.HandleFunc("/", nodeInfoRouter).Methods("GET", "POST")
	myRouter.HandleFunc("/{databaseName}", existsDatabaseRouter).Methods("HEAD")
	myRouter.HandleFunc("/{databaseName}", getDatabaseRouter).Methods("GET")
	myRouter.HandleFunc("/{databaseName}", createNewDatabaseRouter).Methods("PUT")
	myRouter.HandleFunc("/{databaseName}/_all_docs", getAllDocumentsRouter).Methods("GET", "POST")
	myRouter.HandleFunc("/{databaseName}/_bulk_docs", createBulkDocumentsRouter).Methods("POST")
	myRouter.HandleFunc("/{databaseName}/_changes", getDocumentChangesRouter).Methods("GET")
	myRouter.HandleFunc("/{databaseName}/_revs_diff", getRevDiffRouter).Methods("POST")
	myRouter.HandleFunc("/{databaseName}/_ensure_full_commit", ensureCommitRouter).Methods("POST")
	myRouter.HandleFunc("/{databaseName}/_local/{replID}", getReplCPointRecordRouter).Methods("GET")
	myRouter.HandleFunc("/{databaseName}/_local", replCPointRecordRouter).Methods("PUT")
	myRouter.HandleFunc("/{databaseName}/replicate", enableReplicationRouter).Methods("POST")
	myRouter.HandleFunc("/{databaseName}/replicate", disableReplicationRouter).Methods("DELETE")

	// Run server
	fmt.Println("Aquila Port running at localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}
