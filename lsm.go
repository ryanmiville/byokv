package main

import (
	"encoding/json"
	"os"
)

type LSM struct {
	manifest *Manifest
	memtable *MemTable
}

func (lsm *LSM) Get(key string) (string, bool) {
	v, ok := lsm.memtable.Get(key)
	if ok {
		return v, ok
	}
	return lsm.get(key)
}

func (lsm *LSM) Put(key, value string) error {
	if len(lsm.memtable.store) >= 2000 {
		err := lsm.manifest.Flush(lsm.memtable)
		if err != nil {
			return err
		}
	}
	lsm.memtable.Put(key, value)
	return nil
}

func (lsm *LSM) get(key string) (string, bool) {
	cursor := len(lsm.manifest.sstables) - 1
	for i := cursor; i >= 0; i-- {
		filename := lsm.manifest.sstables[i]
		v, ok := scanSSTable(filename, key)
		if ok {
			return v, ok
		}
	}
	return "", false
}

// scanSSTable parses the json file and scans the contents for the key.
func scanSSTable(filename string, key string) (string, bool) {
	// Read the SSTable JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", false
	}

	// Unmarshal into slice of KV structs
	var kvs []KV
	err = json.Unmarshal(data, &kvs)
	if err != nil {
		return "", false
	}

	// Search for the key
	for _, kv := range kvs {
		if kv.Key == key {
			return kv.Value, true
		}
	}

	return "", false
}
