package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"playground/my-dist-kv-store/kvstore"
	"playground/my-dist-kv-store/log"
)

const (
	LogFileName = "transactions.log"
)

var (
	store  *kvstore.KVStore
	logger log.TransactionLogger
)

func init() {
	store = kvstore.NewKVStore()
	err := initTransactionLog()
	if err != nil {
		panic(err)
	}
}

func initTransactionLog() error {
	var err error
	logger, err = log.NewFileTransactionLogger(LogFileName)
	if err != nil {
		return fmt.Errorf("error init transactions log: %w", err)
	}
	events, errs := logger.ReadEvents()
	e, ok := log.Event{}, true

	for err == nil && ok {
		select {
		case err, ok = <-errs:
		case e = <-events:
			switch e.EventType {
			case log.PutEvent:
				err = store.Put(e.Key, e.Value)
			case log.DeleteEvent:
				err = store.Delete(e.Key)
			}
		}
	}

	logger.Run()
	return err
}

// KeyValuePutHandler expects to be called with a PUT request for the "/v1/{Key}" resource.
func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Key"]
	val, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WritePut(key, string(val))
	err = store.Put(key, string(val))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

// KeyValueGetHandler expects to be called with a GET request for the "/v1/{Key}" resource.
func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Key := vars["Key"]

	val, err := store.Get(Key)
	if errors.Is(err, kvstore.ErrKeyNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(val))
}

// KeyValueDeleteHandler expects to be called with a DELETE request for the "/v1/{Key}" resource.
func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Key := vars["Key"]

	logger.WriteDelete(Key)
	err := store.Delete(Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
