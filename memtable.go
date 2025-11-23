package main

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

// MemTable allows thread-safe access to a map.
type MemTable struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewMemTable initializes the store map.
func NewMemTable() *MemTable {
	return &MemTable{
		store: make(map[string]string),
	}
}

// Get retrieves a value safely.
func (s *MemTable) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	return val, ok
}

// Put stores a value safely.
func (s *MemTable) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
}

func (s *MemTable) Flush(filename string) error {
	bytes, err := SortedMapToJSON(s.store)
	if err != nil {
		return err
	}
	s.store = make(map[string]string)
	return os.WriteFile(filename, bytes, 0644)
}

// KV represents the desired JSON object structure inside the array
type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SortedMapToJSON(data map[string]string) ([]byte, error) {
	// 1. Efficiency: Pre-allocate the keys slice to the map size
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}

	// 2. Sort the keys (O(N log N))
	sort.Strings(keys)

	// 3. Efficiency: Pre-allocate the output slice
	// This avoids resizing the underlying array as we add items
	output := make([]KV, len(data))

	for i, k := range keys {
		output[i] = KV{
			Key:   k,
			Value: data[k],
		}
	}

	// 4. Serialize to JSON
	return json.Marshal(output)
}
