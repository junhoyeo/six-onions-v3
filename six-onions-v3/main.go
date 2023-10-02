package main

func main() {
	db := NewInMemoryDB()

	go dnsHandler(db)
	go tcpRequestHandler(db)

	select {} // Block forever to keep the program running
}
