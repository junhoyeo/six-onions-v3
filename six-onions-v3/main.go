package main

func main() {
	var cfg Config
	cfg.Init() // parse flags here

	db := NewAddressMappingTable()

	go dnsHandler(&cfg, db)
	go tcpRequestHandler(&cfg, db)

	select {} // Block forever to keep the program running
}
