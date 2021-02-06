package main

import (
	"fmt"
	"time"
)

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
			fmt.Println("Replication ID generation failed. Replication stopped. ")
			break
		}

		// 4. locate changed documents
		var documents []Document
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
