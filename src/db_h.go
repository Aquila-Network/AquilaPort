package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson"
)

// DBase is Database struct
type DBase struct {
	documentDB *leveldb.DB
	changeDB   *leveldb.DB
	logDB      *leveldb.DB
}

func (db *DBase) createNewDatabase(databaseName string) bool {
	status := true
	dbLocation := DBRoot + "/" + databaseName
	if db.documentDB, err = leveldb.OpenFile(dbLocation+"/documentDB", nil); err != nil {
		fmt.Println(err)
		status = false
	} else if db.changeDB, err = leveldb.OpenFile(dbLocation+"/changeDB", nil); err != nil {
		fmt.Println(err)
		status = false
	} else if db.logDB, err = leveldb.OpenFile(dbLocation+"/logDB", nil); err != nil {
		fmt.Println(err)
		status = false
	}
	return status
}

func (db *DBase) getDocuments(selector string) []Document {
	var documents []Document

	// if selector == "all" {
	// 	// iterate over leveldb and get key, val
	// 	iter := sourceDb.NewIterator(nil, nil)
	// 	for iter.Next() {
	// 		// key := iter.Key()
	// 		value := iter.Value()

	// 		var docRet Document
	// 		// convert bson to byte
	// 		bson.Unmarshal(value, &docRet)

	// 		documents = append(documents, docRet)
	// 	}
	// 	iter.Release()
	// 	// err := iter.Error()
	// }

	return documents
}

func (db *DBase) createNewDocument(w http.ResponseWriter, r *http.Request) {
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

	// // write batch to level db
	// err := sourceDb.Write(batch, nil)
	// fmt.Println(err)

	// // iterate over leveldb and get key, val
	// iter := sourceDb.NewIterator(nil, nil)
	// for iter.Next() {
	// 	key := iter.Key()
	// 	value := iter.Value()

	// 	var docRet Document
	// 	// convert bson to byte
	// 	bson.Unmarshal(value, &docRet)

	// 	fmt.Println(string(key), docRet)
	// }
	// iter.Release()
	// err = iter.Error()

	json.NewEncoder(w).Encode(documents)
}

func (db *DBase) deleteDocument(w http.ResponseWriter, r *http.Request) {
	var ids []string

	json.NewDecoder(r.Body).Decode(&ids)

	// init a batch insert to level
	batch := new(leveldb.Batch)

	// delete docs batch
	for _, id := range ids {
		// delete doc
		batch.Delete([]byte(id))
	}

	// // write batch to level db
	// err := sourceDb.Write(batch, nil)
	// fmt.Println(err)

	// // iterate over leveldb and get key, val
	// iter := sourceDb.NewIterator(nil, nil)
	// for iter.Next() {
	// 	key := iter.Key()
	// 	value := iter.Value()

	// 	var docRet Document
	// 	// convert bson to byte
	// 	bson.Unmarshal(value, &docRet)

	// 	fmt.Println(string(key), docRet)
	// }
	// iter.Release()
	// err = iter.Error()

	json.NewEncoder(w).Encode(ids)
}
