package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson"
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

	existsDatabase(databaseName)

	json.NewEncoder(w).Encode(map[string]string{})
}

func getDatabaseRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	createNewDatabase(databaseName)

	json.NewEncoder(w).Encode(map[string]string{
		"db_name": databaseName,
	})
}

func createNewDatabaseRouter(w http.ResponseWriter, r *http.Request) {
	// get URL params
	params := mux.Vars(r)
	databaseName := params["databaseName"]

	createNewDatabase(databaseName)

	json.NewEncoder(w).Encode(map[string]bool{
		"ok": true,
	})
}

func createNewDocumentRouter(w http.ResponseWriter, r *http.Request) {
	// decode json body
	var documents []Document
	json.NewDecoder(r.Body).Decode(&documents)

	// init a batch insert to level
	batch := new(leveldb.Batch)

	// insert docs to batch
	for _, doc := range documents {
		// update document version
		doc.Version = string(getVersion(doc))
		// convert struct to bson
		data, err := bson.Marshal(doc)
		fmt.Println(err)
		// insert doc
		batch.Put([]byte(doc.ID), data)
	}

	// write batch to level db
	err := sourceDb.Write(batch, nil)
	fmt.Println(err)

	// iterate over leveldb and get key, val
	iter := sourceDb.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var docRet Document
		// convert bson to byte
		bson.Unmarshal(value, &docRet)

		fmt.Println(string(key), docRet)
	}
	iter.Release()
	err = iter.Error()

	json.NewEncoder(w).Encode(documents)
}

func deleteDocumentRouter(w http.ResponseWriter, r *http.Request) {
	var ids []string

	json.NewDecoder(r.Body).Decode(&ids)

	// init a batch insert to level
	batch := new(leveldb.Batch)

	// delete docs batch
	for _, id := range ids {
		// delete doc
		batch.Delete([]byte(id))
	}

	// write batch to level db
	err := sourceDb.Write(batch, nil)
	fmt.Println(err)

	// iterate over leveldb and get key, val
	iter := sourceDb.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		var docRet Document
		// convert bson to byte
		bson.Unmarshal(value, &docRet)

		fmt.Println(string(key), docRet)
	}
	iter.Release()
	err = iter.Error()

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
