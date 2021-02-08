package main

import (
	"encoding/json"
	"fmt"
	"net/http/cookiejar"
)

var jar, err = cookiejar.New(nil)

func authenticate(targetNode string) (int, []byte) {
	return request(targetNode+"/_session", "POST", "name=admin&password=password", "x-www-form-urlencoded")

}

func getNodeInfo(target string) (int, []byte) {
	return request(target, "GET", "", "")
}

func checkDB(root string, dbName string) (int, []byte) {
	return request(root+"/"+dbName, "HEAD", "", "")
}

func createDB(root string, dbName string) (int, []byte) {
	return request(root+"/"+dbName, "PUT", "", "")
}

func getDBInfo(root string, dbName string) (int, []byte) {
	return request(root+"/"+dbName, "GET", "", "")
}

func getReplicationLog(root string, dbName string, logID string) (int, []byte) {
	return request(root+"/"+dbName+"/_local/"+logID, "GET", "", "")
}

func addBatchDocs(root string, dbName string, documents []Document) (int, []byte) {
	data, err := json.Marshal(documents)
	if err != nil {
		fmt.Println(err)
	}

	dataStr := `{"docs":` + string(data) + `}`
	return request(root+"/"+dbName+"/_bulk_docs", "POST", dataStr, "application/json")
}

func getChanges(root string, dbName string) (int, []byte) {
	return request(root+"/"+dbName+"/_changes", "GET", "", "")
}

func getRevsDiffResp(root string, dbName string, revsDiffDoc map[string][]string) (int, []byte) {
	data, err := json.Marshal(revsDiffDoc)
	if err != nil {
		fmt.Println(err)
	}

	dataStr := string(data)
	fmt.Println(dataStr)
	return request(root+"/"+dbName+"/_revs_diff", "POST", dataStr, "")
}

func ensureCommit(root string, dbName string) (int, []byte) {
	return request(root+"/"+dbName+"/_ensure_full_commit", "POST", "", "application/json")
}

func setReplChkPoint(root string, dbName string, replLog []byte) (int, []byte) {
	return request(root+"/"+dbName+"/_ensure_full_commit", "POST", string(replLog), "application/json")
}
