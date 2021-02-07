package main

import (
	"fmt"
	"regexp"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

var c conf
var databases = make(map[string]DBase)

// remove all non alphanumeric chars from string
var nalReg, _ = regexp.Compile("[^a-zA-Z0-9]+")

// DBRoot is database root directory on disk
var DBRoot = c.getConf(configFile).DBRoot

func existsDatabase(databaseName string) bool {
	if _, ok := databases[databaseName]; ok {
		return true
	}
	return false
}

func createNewDatabase(databaseName string) bool {
	status := true
	dbObj := DBase{}
	status = dbObj.createNewDatabase(databaseName)

	if status {
		databases[databaseName] = dbObj
	}

	return status
}

func getAllDocuments(databaseName string) []Document {
	if dbObj, ok := databases[databaseName]; ok {
		return dbObj.getDocuments("all")
	}
	return nil
}

func createNewDocuments(databaseName string, documents []Document) []Document {

	if dbObj, ok := databases[databaseName]; ok {
		return dbObj.createNewDocuments(documents)
	}
	return nil
}

func getDocumentChanges(databaseName string) ChangeDocument {
	if dbObj, ok := databases[databaseName]; ok {
		return dbObj.getChanges("0", 100)
	}
	return ChangeDocument{}
}

func getRevDiff(databaseName string, changeMap map[string][]string) map[string][]string {
	// check version of existing docs < change version, need to replicate
	diffMap := make(map[string][]string)

	if dbObj, ok := databases[databaseName]; ok {
		for ID, version := range changeMap {
			docIn, err := dbObj.documentDB.Get([]byte(ID), nil)
			if err == nil {
				doc := Document{}
				err := bson.Unmarshal(docIn, &doc)
				if err == nil {
					versionDoc, _ := strconv.Atoi(nalReg.ReplaceAllString(doc.Version, ""))
					versionRev, _ := strconv.Atoi(nalReg.ReplaceAllString(version[0], ""))

					if versionDoc < versionRev {
						// needs to be replicated
						diffMap[ID] = version
					}
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}

	return diffMap
}

func replCPointRecord(databaseName string, rcpoint ReplCheckpoint) bool {
	if dbObj, ok := databases[databaseName]; ok {
		data, err := bson.Marshal(rcpoint)
		if err == nil {
			err := dbObj.logDB.Put([]byte(rcpoint.ID), data, nil)
			if err == nil {
				return true
			}
		}
	}
	return false
}
