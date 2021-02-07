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

func deleteDocument(ids []string) []string {

	return ids
}
