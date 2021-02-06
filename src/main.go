// main.go
package main

func main() {
	var c conf
	port := c.getConf(configFile).Port

	go replicatorDemon()

	handleRequests(port)
}
