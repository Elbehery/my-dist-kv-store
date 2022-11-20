package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"playground/my-dist-kv-store/kvstore"
)

var store kvstore.KVStore

func init() {
	store = kvstore.NewKVStore()
}

// KeyValuePutHandler expects to be called with a PUT request for the "/v1/{key}" resource.
func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = store.Put(key, string(val))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

// KeyValueGetHandler expects to be called with a GET request for the "/v1/{key}" resource.
func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	val, err := store.Get(key)
	if errors.Is(err, kvstore.ErrKeyNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(val))
}

// KeyValueDeleteHandler expects to be called with a DELETE request for the "/v1/{key}" resource.
func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
