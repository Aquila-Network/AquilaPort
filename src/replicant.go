package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func updateReplicationStatus(databaseName string, tID string, replReq ReplRequest, status string) bool {
	if _, ok := databases[databaseName]; ok {
		target := replReq.Target

		// ========= Generate replication ID =================
		// get source node ID
		sNodeID := getPortUUID()

		if tID == "" {
			// get target node ID
			rStatus, tNodeIDB := getNodeInfo(target)

			if rStatus == 200 {
				var tNodeID = NodeStatus{}
				if err := json.Unmarshal(tNodeIDB, &tNodeID); err != nil {
					fmt.Println(err)
				} else {
					tID = tNodeID.NodeID
				}
			}
		}

		if tID != "" {
			sha := getReplID(sNodeID, tID, databaseName)

			// Add replication ID with details to replication DB
			replData := map[string]string{}
			replData["ID"] = sha
			replData["source"] = sNodeID
			replData["target"] = tID
			replData["targetUrl"] = target
			replData["status"] = status
			replData["databaseName"] = databaseName

			data, err := bson.Marshal(replData)
			fmt.Println(sha)
			if err == nil {
				replicationDB.Put([]byte(sha), data, nil)
				return true
			} else {
				fmt.Println(err)
			}
		}

	}
	return false
}

func replicatorDaemon() {
	// wait 5 seconds before starting daemon
	time.Sleep(time.Duration(1 * time.Second))

	fmt.Println("Starting replicant, go..")

	for {
		iter := replicationDB.NewIterator(nil, nil)
		for iter.Next() {
			replicationID := string(iter.Key())
			value := iter.Value()
			replData := map[string]string{}
			err := bson.Unmarshal(value, &replData)

			if err == nil && replData["status"] == "active" {
				// initial data
				u, _ := url.Parse(replData["targetUrl"])
				targetRoot := "http://" + u.Host
				targetNodeIsAlive := true
				// Set databases
				targetDB := replData["databaseName"]

				// validate node id and availability
				sNodeID := getPortUUID()
				rStatus, tNodeIDB := getNodeInfo(targetRoot)
				var tNodeID = NodeStatus{}

				if rStatus == 200 {
					if err := json.Unmarshal(tNodeIDB, &tNodeID); err != nil {
						fmt.Println(err)
					}
					sha := getReplID(sNodeID, tNodeID.NodeID, targetDB)
					if replicationID != sha {
						// target Aquila Port node is not found
						targetNodeIsAlive = false
					}
				} else {
					// target Aquila Port node is not alive
					targetNodeIsAlive = false
				}

				// Pause replication for dead target node
				if !targetNodeIsAlive {
					fmt.Println("Target node " + replData["target"] + " is not alive. Disabling replication (ID: " + replicationID + ") for target.")
					replReq := ReplRequest{}
					replReq.Target = targetRoot
					updateReplicationStatus(targetDB, replData["target"], replReq, "inactive")
					continue
				}

				// TODO: authenticate couchDB

				// 1. verify peers
				stat, _ := checkDB(targetRoot, targetDB)
				if stat == 404 {
					fmt.Println("Target DB do not exist. Creating it..")
					// create target database
					stat, _ = createDB(targetRoot, targetDB)
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
				stat, _ = getDBInfo(targetRoot, targetDB)
				if stat == 200 {
					// fmt.Println("-")
				}

				// 3. find common ascestry
				fullReplication := false
				var replicationLog ReplCheckpoint

				if replicationID != "" {
					fmt.Println("Generated replication ID: ", replicationID)
					// get replication log from target
					fmt.Println("Getting replication log from target..")
					stat, rlog := getReplicationLog(targetRoot, targetDB, replicationID)
					if stat == 200 {
						// compare logs
						if err := json.Unmarshal(rlog, &replicationLog); err != nil {
							fmt.Println(err)
						}
					} else if stat == 404 {
						fullReplication = true
						fmt.Println("Replication log not available. Full replication needed.")
					} else {
						fullReplication = true
						fmt.Println("Unknown log format. Full replication needed.")
					}
				} else {
					fmt.Println("Replication ID generation failed. Replication stopped. ")
					continue
				}

				// 4. locate changed documents
				var changes ChangeDocument
				var replLoguntil string
				if fullReplication {
					// get all changes
					replLoguntil, changes = getDocumentChanges(targetDB, "0")
				} else {
					since := replicationLog.Rev
					// Get changed documents in target
					replLoguntil, changes = getDocumentChanges(targetDB, since)
					fmt.Println(changes)
				}
				// 4.1 prepare revs diff docuemnt
				revsDiffDoc := map[string][]string{}
				changeCount := 0

				for _, element := range changes.Results {
					docID := element.ID
					rev := []string{""}
					rev[0] = element.Changes[0].Rev

					revsDiffDoc[docID] = rev
					changeCount++
				}
				// skip replication step if no changes found
				if changeCount <= 0 {
					fmt.Println("No changes, Skip replication.")
					continue
				}
				// 4.2 get revs diff confirmation
				var wantReplDocs map[string][]string

				stat, wantRepl := getRevsDiffResp(targetRoot, targetDB, revsDiffDoc)
				if stat == 200 {
					// compare logs
					if err := json.Unmarshal(wantRepl, &wantReplDocs); err != nil {
						fmt.Println(err)
					}
				} else if stat == 404 {
					fmt.Println("Revs. diff. document retrieval failed. Skipping replication.")
					continue
				}

				// 5. replicate changes
				documentIds := []string{}
				changeCount = 0
				for docID, version := range wantReplDocs {
					documentIds = append(documentIds, docID)
					if version[0] == "" {
						fmt.Println("invalid version")
					} else {
						changeCount++
					}
				}
				if changeCount <= 0 {
					fmt.Println("No rev. changes.")
				}

				documents := getBulkDocuments(targetDB, documentIds)
				if len(documents) > 0 {
					stat, _ := addBatchDocs(targetRoot, targetDB, documents)
					fmt.Println(stat)
					if stat == 201 {
						fmt.Println("Documents written succesfully.")
					} else {
						fmt.Println("Documents write failed.")
					}
				}

				// ensure in commit
				fmt.Println(replLoguntil)
				stat, data := ensureCommit(targetRoot, targetDB)
				if stat == 201 {
					fmt.Println("Documents commited succesfully.")

				} else {
					fmt.Println("Documents commit failed.", string(data))
					break
				}

				// set record replication checkpoint
				replCPoint := ReplCheckpoint{}
				replCPoint.ID = replicationID
				replCPoint.Ok = true
				replCPoint.Rev = replLoguntil
				stat, _ = setReplChkPoint(targetRoot, targetDB, replCPoint)
			}

			// wait 5 seconds before next replication loop
			time.Sleep(time.Duration(5 * time.Second))

		}
		iter.Release()
		err = iter.Error()

		// wait 5 seconds before next replication loop
		time.Sleep(time.Duration(5 * time.Second))
	}
}
