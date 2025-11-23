package main

import "net/http"

func main() {
	manifest, err := GetOrCreateManifest()
	if err != nil {
		panic(err)
	}
	memtable := NewMemTable()

	lsm := LSM{manifest: manifest, memtable: memtable}

	server := Server(lsm)
	http.ListenAndServe(":8080", server)
}
