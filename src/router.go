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

func createNewDocument(w http.ResponseWriter, r *http.Request) {
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

func deleteDocument(w http.ResponseWriter, r *http.Request) {
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
	myRouter.HandleFunc("/create", createNewDocument).Methods("POST")
	myRouter.HandleFunc("/delete", deleteDocument).Methods("POST")

	// Run server
	fmt.Println("Aquila Port running at localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}
