package main

func main() {
	db := NewAddressMappingTable()

	go dnsHandler(db)
	go tcpRequestHandler(db)

	select {} // Block forever to keep the program running
}
