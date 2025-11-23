package main

import (
	"fmt"
	"io"
	"net/http"
)

func Server(lsm LSM) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		value, exists := lsm.Get(key)

		if !exists {
			http.Error(w, fmt.Sprintf("key not found: %s", key), http.StatusNotFound)
			return
		}

		fmt.Fprint(w, value)
	})

	mux.HandleFunc("POST /{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		value := string(body)

		lsm.Put(key, value)

		fmt.Fprintf(w, "stored value for key: %s", key)
	})

	return mux
}
