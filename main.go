package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"log"
	"net/http"
)

type cacheObjectStruct struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

var host = flag.String("host", ":8080", "host to bind to")
var cert = flag.String("cert", "", "tls certificate")
var key = flag.String("key", "", "tls certificate key")
var useTLS = flag.Bool("tls", false, "enable TLS")
var c = cache.New(cache.NoExpiration, 0)

func GetCacheHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	v, found := c.Get(id)
	if found {
		res := &cacheObjectStruct{ID: id, Value: v.(string)}
		s, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, string(s))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func PostCacheHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var cacheObject cacheObjectStruct
	err = json.Unmarshal(body, &cacheObject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		c.Add(cacheObject.ID, cacheObject.Value, cache.NoExpiration)
		w.WriteHeader(http.StatusOK)
	}

}

func DeleteCacheHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	_, found := c.Get(id)
	if found {
		c.Delete(id)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func PutCacheHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var cacheObject cacheObjectStruct
		err = json.Unmarshal(body, &cacheObject)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			_, found := c.Get(id)
			if found {
				c.Delete(id)
				c.Add(cacheObject.ID, cacheObject.Value, cache.NoExpiration)
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func main() {

	flag.Parse()

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/cache", PostCacheHandler).Methods("POST")
	r.HandleFunc("/cache/{id}", GetCacheHandler).Methods("GET")
	r.HandleFunc("/cache/{id}", DeleteCacheHandler).Methods("DELETE")
	r.HandleFunc("/cache/{id}", PutCacheHandler).Methods("PUT")

	// HTTP
	if !*useTLS {
		fmt.Printf("listening on %s /cache \r\n", *host)
		err := http.ListenAndServe(*host, r)
		if err != nil {
			log.Printf("error : %v", err)
		}
	}

	// HTTPS
	if *useTLS {
		fmt.Printf("listening on %s /cache cert=%s, key=%s\r\n", *host, *cert, *key)
		err := http.ListenAndServeTLS(*host, *cert, *key, r)
		if err != nil {
			log.Printf("error : %v", err)
		}
	}

}
