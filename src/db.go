package main

var c conf
var databases = make(map[string]DBase)

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

func getDocuments(selector string) []Document {
	var documents []Document

	return documents
}

func createNewDocument(documents []Document) []Document {

	return documents
}

func deleteDocument(ids []string) []string {

	return ids
}
