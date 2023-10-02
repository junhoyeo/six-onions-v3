package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func tcpRequestHandler(db *InMemoryDB) {
	tport := flag.Int("transport", 1337,
		"the port that iptables will be redirecting connections to")
	flag.Parse()

	la, _ := net.ResolveTCPAddr("tcp6", fmt.Sprintf("[::]:%d", *tport))
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

		go handleConn(c)
	}
}
