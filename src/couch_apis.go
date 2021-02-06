package main

import (
	"encoding/json"
	"fmt"
	"net/http/cookiejar"
)

var jar, err = cookiejar.New(nil)

func authenticate() (int, []byte) {
	return request("http://127.0.0.1:5984/_session", "POST", "name=admin&password=password", "x-www-form-urlencoded")

}

func checkDB(dbName string) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName, "HEAD", "", "")
}

func createDB(dbName string) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName, "PUT", "", "")
}

func getDBInfo(dbName string) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName, "GET", "", "")
}

func getReplicationLog(dbName string, logID string) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName+"/_local/"+logID, "GET", "", "")
}

func addBatchDocs(dbName string, documents []Document) (int, []byte) {
	data, err := json.Marshal(documents)
	if err != nil {
		fmt.Println(err)
	}

	dataStr := `{"docs":` + string(data) + `}`
	return request("http://127.0.0.1:5984/"+dbName+"/_bulk_docs", "POST", dataStr, "application/json")
}

func getChanges(dbName string) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName+"/_changes?style=all_docs", "GET", "", "")
}

func ensureCommit(dbName string) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName+"/_ensure_full_commit", "POST", "", "application/json")
}

func setReplChkPoint(dbName string, replLog []byte) (int, []byte) {
	return request("http://127.0.0.1:5984/"+dbName+"/_ensure_full_commit", "POST", string(replLog), "application/json")
}
