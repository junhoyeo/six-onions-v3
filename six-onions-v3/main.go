package main

import (
	"encoding/base32"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/net/proxy"
)

func main() {
	db := NewInMemoryDB()

	go dnsHandler(db)
	go tcpRequestHandler(db)

	select {} // Block forever to keep the program running
}

func handleConn(c *net.TCPConn) {
	// first, let's recover the address
	tc, fd, err := realServerAddress(c)
	defer c.Close()
	defer fd.Close()

	if err != nil {
		log.Printf("Unable to recover address %s", err.Error())
		return
	}

	toraddr := tc.IP[6:]
	toronionaddr :=
		fmt.Sprintf("%s.onion", base32.StdEncoding.EncodeToString(toraddr))

	if !isAllowedPort(tc.Port) {
		log.Printf("Disallowed connection from %s to %s:%d due to port block",
			c.RemoteAddr().String(), toronionaddr, tc.Port)
		return
	}

	log.Printf("Connection from %s to %s:%d",
		c.RemoteAddr().String(), toronionaddr, tc.Port)

	d, err := proxy.SOCKS5("tcp", "localhost:9050", nil, proxy.Direct)
	if err != nil {
		log.Printf("Unable to recover address %s", err.Error())
		return
	}

	torconn, err := d.Dial("tcp", fmt.Sprintf("%s:%d", toronionaddr, tc.Port))
	if err != nil {
		log.Printf("Tor conncetion error %s", err.Error())
		return
	}

	go io.Copy(torconn, fd)
	io.Copy(fd, torconn)
}
