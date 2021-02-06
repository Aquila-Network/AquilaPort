// main.go
package main

func main() {
	go replicatorDemon()

	handleRequests("5006")
}
