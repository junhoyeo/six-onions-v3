package main

import (
	"encoding/binary"
	"net"
	"sync"
)

// AddressMappingTable holds the mappings between OnionV3 addresses and IPv6 addresses.
type AddressMappingTable struct {
	data     map[string]string // IPv6 to OnionV3 mapping
	mu       sync.RWMutex
	lastIPv6 [16]byte
}

// NewAddressMappingTable initializes a new AddressMappingTable instance.
func NewAddressMappingTable() *AddressMappingTable {
	// Set initial address
	initialIPv6 := [16]byte{0x2a, 0x0c, 0x2f, 0x07, 0xFE, 0xD5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	return &AddressMappingTable{
		data:     make(map[string]string),
		lastIPv6: initialIPv6,
	}
}

// Set adds a mapping between the provided IPv6 and OnionV3 addresses.
func (db *AddressMappingTable) Set(IPv6, OnionV3 string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[IPv6] = OnionV3
}

// NextIPv6 generates the next IPv6 address sequentially.
func (db *AddressMappingTable) NextIPv6() net.IP {
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
func (db *AddressMappingTable) GetByIPv6(IPv6 string) (OnionV3 string, ok bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	OnionV3, ok = db.data[IPv6]
	return
}
