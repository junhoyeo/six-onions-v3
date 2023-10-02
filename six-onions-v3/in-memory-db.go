package main

import (
	"encoding/binary"
	"net"
	"sync"
)

type InMemoryDB struct {
	data        map[string]string // OnionV3Address to IPv6 mapping
	reverseData map[string]string // IPv6 to OnionV3Address mapping
	mu          sync.RWMutex
	lastAddress [16]byte
}

func NewInMemoryDB() *InMemoryDB {
	// Set initial address
	initialAddress := [16]byte{0x2a, 0x0c, 0x2f, 0x07, 0xFE, 0xD5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	return &InMemoryDB{
		data:        make(map[string]string),
		reverseData: make(map[string]string),
		lastAddress: initialAddress,
	}
}

func (db *InMemoryDB) Set(OnionV3Address, IPv6 string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[OnionV3Address] = IPv6
	db.reverseData[IPv6] = OnionV3Address
}

func (db *InMemoryDB) GetByOnionV3Address(OnionV3Address string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	IPv6, ok := db.data[OnionV3Address]
	return IPv6, ok
}

func (db *InMemoryDB) NextIPv6() net.IP {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Increment the last 64 bits of the address to get the next address
	// (assuming you only want to increment the host portion of the address)
	nextValue := binary.BigEndian.Uint64(db.lastAddress[8:]) + 1
	binary.BigEndian.PutUint64(db.lastAddress[8:], nextValue)

	// Copy the last address to ensure the underlying array isn't modified outside this method
	nextAddress := make(net.IP, net.IPv6len)
	copy(nextAddress, db.lastAddress[:])

	return nextAddress
}

func (db *InMemoryDB) GetByIPv6(IPv6 string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	OnionV3Address, ok := db.reverseData[IPv6]
	return OnionV3Address, ok
}
