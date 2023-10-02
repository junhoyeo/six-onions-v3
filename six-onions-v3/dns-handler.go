package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func dnsHandler(cfg *Config, db *AddressMappingTable) {
	log.Printf("listening on %s for %s", cfg.DnsListen, cfg.Domain)

	// how do we override the dnsbase?
	// ./six-onions-v3 -domain=tor6.flm.me.uk -listen=

	listens, err := net.ListenPacket("udp4", cfg.DnsListen)
	if err != nil {
		log.Fatalf("failed to listen on UDP %s / %s", cfg.DnsListen, err.Error())
	}

	for {
		dnsin := make([]byte, 1500)
		inbytes, inaddr, err := listens.ReadFrom(dnsin)

		inmsg := &dns.Msg{}

		if unpackErr := inmsg.Unpack(dnsin[0:inbytes]); unpackErr != nil {
			log.Printf("Unable to unpack DNS request %s", err.Error())
			continue
		}

		if len(inmsg.Question) != 1 {
			log.Printf("More than one quesion in query (%d), droppin %+v", len(inmsg.Question), inmsg)
			continue
		}

		iqn := strings.ToLower(inmsg.Question[0].Name)

		if !strings.Contains(iqn, cfg.Domain) {
			log.Printf("question is not for us '%s' vs expected '%s'", iqn, cfg.Domain)
			continue
		}

		outmsg := &dns.Msg{}

		iqn = strings.ToUpper(inmsg.Question[0].Name)

		queryname := strings.Replace(
			iqn, fmt.Sprintf(".%s.", strings.ToUpper(cfg.Domain)), "", 1)

		// Generate a new IPv6 address
		newIPv6 := db.NextIPv6()
		// Convert net.IP to string to store in the DB
		newIPv6Str := newIPv6.String()
		// Store in the database
		db.Set(queryname, newIPv6Str)

		outmsg.Id = inmsg.Id
		outmsg = inmsg.SetReply(outmsg)
		iqn = inmsg.Question[0].Name
		outmsg.Answer = make([]dns.RR, 1)
		outmsg.Answer[0] = &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   iqn,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    2147483646},
			AAAA: newIPv6,
		}
		outputb, err := outmsg.Pack()

		if err != nil {
			log.Printf("unable to pack response to thing %s", err.Error())
			continue
		}

		listens.WriteTo(outputb, inaddr)
	}
}
