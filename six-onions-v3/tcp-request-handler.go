package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

func tcpRequestHandler(cfg *Config, db *AddressMappingTable) {
	la, _ := net.ResolveTCPAddr("tcp6", fmt.Sprintf("[::]:%d", cfg.Transport))
	l, err := net.ListenTCP("tcp6", la)
	if err != nil {
		log.Fatalf("Unable to listen on the transparent port %s",
			err.Error())
	}

	failurecount := 0
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			if failurecount != 50 {
				failurecount++
			} else {
				log.Printf("Unable to accept connection! %s", err.Error())
			}
			time.Sleep(time.Millisecond * time.Duration(failurecount*10))
			continue
		}
		failurecount = 0

		go handleTCPConn(c, db)
	}
}

func handleTCPConn(c *net.TCPConn, db *AddressMappingTable) {
	// first, let's recover the address
	tc, fd, err := realServerAddress(c)
	defer c.Close()
	defer fd.Close()

	if err != nil {
		log.Printf("Unable to recover address %s", err.Error())
		return
	}

	// Convert the IP part to string, as your DB uses string keys and values
	ipStr := tc.IP.String()
	// Lookup the OnionV3Address from the DB using the IP address as the key
	onionV3Address, ok := db.GetByIPv6(ipStr)
	if !ok {
		log.Printf("No OnionV3Address found for IP: %s", ipStr)
		return
	}
	toronionaddr := fmt.Sprintf("%s.onion", onionV3Address)

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
