package main

import (
	"flag"
)

type Config struct {
	Domain    string
	DnsListen string
	Transport int
}

func (cfg *Config) Init() {
	flag.StringVar(&cfg.Domain, "domain", "tor6.flm.me.uk", "the domain you want to top on")
	flag.StringVar(&cfg.DnsListen, "listen", "127.0.0.1:553", "the port to listen on")
	flag.IntVar(&cfg.Transport, "transport", 1337, "the port that iptables will be redirecting connections to")
	flag.Parse()
}
