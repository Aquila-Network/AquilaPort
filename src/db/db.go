package db

import (
	"/util/ctypes"

	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	_localDb, _      = leveldb.OpenFile("./db/_local", nil)
	replicationDb, _ = leveldb.OpenFile("./db/replication", nil)
	sourceDb, _      = leveldb.OpenFile("./db/source", nil)
)

func getDocuments(selector string) []ctypes.Document {
	var documents []ctypes.Document

	if selector == "all" {
		// iterate over leveldb and get key, val
		iter := sourceDb.NewIterator(nil, nil)
		for iter.Next() {
			// key := iter.Key()
			value := iter.Value()

			var docRet ctypes.Document
			// convert bson to byte
			bson.Unmarshal(value, &docRet)

			documents = append(documents, docRet)
		}
		iter.Release()
		// err := iter.Error()
	}

	return documents
}
