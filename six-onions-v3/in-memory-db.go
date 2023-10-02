package main

import (
	"encoding/binary"
	"net"
	"sync"
)

// InMemoryDB holds the mappings between OnionV3 addresses and IPv6 addresses.
type InMemoryDB struct {
	data        map[string]string // OnionV3 to IPv6 mapping
	reverseData map[string]string // IPv6 to OnionV3 mapping
	mu          sync.RWMutex
	lastIPv6    [16]byte
}

// NewInMemoryDB initializes a new InMemoryDB instance.
func NewInMemoryDB() *InMemoryDB {
	// Set initial address
	initialIPv6 := [16]byte{0x2a, 0x0c, 0x2f, 0x07, 0xFE, 0xD5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	return &InMemoryDB{
		data:        make(map[string]string),
		reverseData: make(map[string]string),
		lastIPv6:    initialIPv6,
	}
}

// Set adds a mapping between the provided OnionV3 and IPv6 addresses.
func (db *InMemoryDB) Set(OnionV3, IPv6 string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[OnionV3] = IPv6
	db.reverseData[IPv6] = OnionV3
}

// NextIPv6 generates the next IPv6 address sequentially.
func (db *InMemoryDB) NextIPv6() net.IP {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Increment the last 64 bits of the address to get the next address.
	// (assuming you only want to increment the host portion of the address)
	nextValue := binary.BigEndian.Uint64(db.lastIPv6[8:]) + 1
	binary.BigEndian.PutUint64(db.lastIPv6[8:], nextValue)

	// Copy the last address to ensure the underlying array isn't modified outside this method.
	nextIPv6 := make(net.IP, net.IPv6len)
	copy(nextIPv6, db.lastIPv6[:])

	return nextIPv6
}

// GetByIPv6 retrieves the OnionV3 address mapped to the provided IPv6 address.
func (db *InMemoryDB) GetByIPv6(IPv6 string) (OnionV3 string, ok bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	OnionV3, ok = db.reverseData[IPv6]
	return
}
