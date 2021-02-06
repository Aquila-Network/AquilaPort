// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/Aquila-Network/AquilaPort/src/util/ctypes"
)

var jar, err = cookiejar.New(nil)

// if err != nil {
//     fmt.Println(err)
// }

// =========================== COUCHDB ======================================================================
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

func addBatchDocs(dbName string, documents []ctypes.Document) (int, []byte) {
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

func request(url string, method string, payload string, contentType string) (int, []byte) {

	client := &http.Client{
		Jar: jar,
	}

	var req *http.Request
	var err error

	if payload == "" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(payload))
	}

	if err != nil {
		fmt.Println(err)
	}

	if contentType == "x-www-form-urlencoded" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else if contentType == "application/json" {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return res.StatusCode, body
}

func replicatorDemon() {
	for {
		// authenticate couchDB
		stat, _ := authenticate()
		if stat != 200 {
			fmt.Println("CouchDB Authentication failed. Exiting demon.")
		} else {
			fmt.Println("CouchDB authentication success.")
		}

		// perform replication to target =======================================
		// sourceDB := "source"
		targetDB := "target"

		// 1. verify peers
		stat, _ = checkDB(targetDB)
		if stat == 404 {
			fmt.Println("Target DB do not exist. Creating it..")
			// create target database
			stat, _ = createDB(targetDB)
			if stat == 201 {
				fmt.Println("Target DB created.")
			} else {
				fmt.Println("Target DB can not be created. Aborting replication..")
				break
			}
		} else if stat == 200 {
			fmt.Println("Target DB exists")
		}
		// 2. get peers information
		fmt.Println("Getting peers information..")
		stat, _ = getDBInfo(targetDB)
		replicationID := ""
		if stat == 200 {
			// generate replication ID
			replicationID = "123456"
		}

		// 3. find common ascestry
		fullReplication := false

		if replicationID != "" {
			fmt.Println("Generated replication ID: ", replicationID)
			// get replication log from target
			fmt.Println("Getting replication log from target..")
			stat, rlog := getReplicationLog(targetDB, replicationID)
			if stat == 200 {
				// compare logs
				fmt.Println(string(rlog))
			} else if stat == 404 {
				fullReplication = true
				fmt.Println("Replication log not available. Full replication needed.")
			}
		} else {
			fmt.Println("Replication ID generation failed. Exiting.. ")
			break
		}

		// 4. locate changed documents
		var documents []ctypes.Document
		if fullReplication {
			// get all documents
			documents = getDocuments("all")
		} else {
			// Get changed documents in target
			stat, changes := getChanges(targetDB)
			if stat == 200 {
				fmt.Println("Changes: ", changes)
			}
			// finalize documents to be replicated
			documents = getDocuments("all") // TODO: To be changed to selectives

			// finalize replication if no change found
			if len(documents) <= 0 {
				fmt.Println("No more changes to replicate.")
			}
		}

		// 5. replicate changes
		if len(documents) > 0 {
			stat, _ := addBatchDocs(targetDB, documents)
			if stat == 201 {
				fmt.Println("Documents written succesfully.")
			} else {
				fmt.Println("Documents write failed.")
			}

			// ensure in commit
			stat, data := ensureCommit(targetDB)
			if stat == 201 {
				fmt.Println("Documents commited succesfully.")

			} else {
				fmt.Println("Documents commit failed.", string(data))
				break
			}

			// set record replication checkpoint
			stat, _ = setReplChkPoint(targetDB, []byte(""))
		}

		// wait 5 seconds before next replication loop
		time.Sleep(time.Duration(5 * time.Second))
	}
}

func main() {
	go replicatorDemon()

	handleRequests()
}
