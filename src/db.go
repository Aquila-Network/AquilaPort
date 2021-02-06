package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	_localDb, _      = leveldb.OpenFile("./db/_local", nil)
	replicationDb, _ = leveldb.OpenFile("./db/replication", nil)
	sourceDb, _      = leveldb.OpenFile("./db/source", nil)
)

func existsDatabase(databaseName string) string {
	return databaseName
}

func createNewDatabase(databaseName string) string {
	return databaseName
}

func getDocuments(selector string) []Document {
	var documents []Document

	if selector == "all" {
		// iterate over leveldb and get key, val
		iter := sourceDb.NewIterator(nil, nil)
		for iter.Next() {
			// key := iter.Key()
			value := iter.Value()

			var docRet Document
			// convert bson to byte
			bson.Unmarshal(value, &docRet)

			documents = append(documents, docRet)
		}
		iter.Release()
		// err := iter.Error()
	}

	return documents
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
