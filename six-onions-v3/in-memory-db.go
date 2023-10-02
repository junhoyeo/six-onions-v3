package main

import (
	"sync"
)

type InMemoryDB struct {
	data        map[string]string // OnionV3Address to IPv6 mapping
	reverseData map[string]string // IPv6 to OnionV3Address mapping
	mu          sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data:        make(map[string]string),
		reverseData: make(map[string]string),
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

func (db *InMemoryDB) GetByIPv6(IPv6 string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	OnionV3Address, ok := db.reverseData[IPv6]
	return OnionV3Address, ok
}
