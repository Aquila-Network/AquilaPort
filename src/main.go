// main.go
package main

func main() {
	var c conf
	port := c.getConf(configFile).Port

	// initialize databases
	initDatabases()
	// set a random Aquila Port node ID

	// start replication daemon
	go replicatorDaemon()
	// start server & handle API calls
	handleRequests(port)
}
