package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson"
)

// DBase is Database struct
type DBase struct {
	documentDB *leveldb.DB
	changeDB   *leveldb.DB
	logDB      *leveldb.DB
}

var localDB *leveldb.DB
var replicationDB *leveldb.DB

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

	// ChangeDB set initial value
	_, err := db.changeDB.Get([]byte("0"), nil)
	if err != nil {
		db.changeDB.Put([]byte("0"), []byte("0"), nil)
	}

	// update localDB
	err = localDB.Put([]byte("LDB_"+databaseName), []byte("active"), nil)
	if err != nil {
		fmt.Println(err)
	}

	return status
}

func (db *DBase) getDocuments(selector string) []Document {
	var documents []Document

	if selector == "all" {
		// iterate over leveldb and get key, val
		iter := db.documentDB.NewIterator(nil, nil)
		for iter.Next() {
			// key := iter.Key()
			value := iter.Value()

			var docRet Document
			// convert bson to byte
			bson.Unmarshal(value, &docRet)

			documents = append(documents, docRet)
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			fmt.Println(err)
		}
	}

	return documents
}

func (db *DBase) getDocumentsByIds(ids []string) []Document {
	var documents []Document

	for _, id := range ids {
		data, err := db.documentDB.Get([]byte(id), nil)
		if err == nil {
			document := Document{}

			err := bson.Unmarshal(data, &document)
			if err == nil {
				documents = append(documents, document)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

	return documents
}

func (db *DBase) createNewDocuments(documents []Document) []Document {
	// init a batch insert to level
	batch := new(leveldb.Batch)
	// append ids
	var ids string
	// array for docs
	var docs []Document

	// insert docs to batch
	for _, doc := range documents {
		// update document version
		doc.Version = string(getVersion(doc))
		// convert struct to bson
		data, err := bson.Marshal(doc)
		if err != nil {
			fmt.Println(err)
		}
		// insert doc
		batch.Put([]byte(doc.ID), data)
		if ids == "" {
			ids = doc.ID
		} else {
			ids = ids + "|" + doc.ID
		}
		docs = append(docs, doc)
	}

	// write batch to level db
	err := db.documentDB.Write(batch, nil)
	if err != nil {
		fmt.Println(err)
	}

	// update changes DB
	db.changeDB.Put([]byte(strconv.FormatInt(time.Now().Unix(), 10)), []byte(ids), nil)

	return docs
}

func (db *DBase) getChanges(since string, max int) (string, ChangeDocument) {
	// var changes []ChangeDocument

	iter := db.changeDB.NewIterator(nil, nil)
	ok := iter.Seek([]byte(since))
	ok = iter.Next()
	results := []ChangeResultsDocument{}
	lastSeq := 0
	untilKey := "0"
	changeMap := map[string]ChangeResultsDocument{}
	for ; ok; ok = iter.Next() {
		// Use key, value
		untilKey = string(iter.Key())
		result := ChangeResultsDocument{}
		for _, element := range strings.Split(string(iter.Value()), "|") {
			// get corresponding document
			docIn, err := db.documentDB.Get([]byte(element), nil)
			if err == nil {
				doc := Document{}
				err := bson.Unmarshal(docIn, &doc)
				if err == nil {
					changes := []ChangeResultsChangesDocument{}
					changes = append(changes, ChangeResultsChangesDocument{doc.Version})
					result.Changes = changes
					result.ID = element
					result.Seq = lastSeq
					result.Deleted = doc.Deleted

					//  update changeMap to eliminate duplication
					changeMap[element] = result
					lastSeq = lastSeq + 1
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
	iter.Release()
	err = iter.Error()

	// build results from map
	// and return changes
	for _, result := range changeMap {
		results = append(results, result)
	}
	return untilKey, ChangeDocument{lastSeq - 1, 0, results}
}

func (db *DBase) commitChanges() bool {
	return true
}

func (db *DBase) getReplLog() {
	// return true
}
