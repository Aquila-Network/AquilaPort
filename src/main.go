// main.go
package main

func main() {
	var c conf
	port := c.getConf("../port_config.yml").Port

	go replicatorDemon()

	handleRequests(port)
}
