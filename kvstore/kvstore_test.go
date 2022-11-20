package kvstore

import (
	"errors"
	"os"
	"testing"
)

var store *KVStore

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestKVStore_Put(t *testing.T) {

	scenarios := []struct {
		name   string
		Key    string
		val    string
		expErr error
	}{
		{"new Key",
			"new Key",
			"new value",
			nil,
		},
		{
			"old Key",
			"old Key",
			"old value",
			nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := store.Put(scenario.Key, scenario.val)
			if scenario.expErr != nil {
				if err == nil {
					t.Errorf("expected error %v, but got %v instead\n", scenario.expErr, err)
				}
				if !errors.Is(err, scenario.expErr) {
					t.Errorf("expected error %v, but got %v instead\n", scenario.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error %v, but got %v instead\n", scenario.expErr, err)
			}
		})
	}
}

func TestKVStore_Get(t *testing.T) {

	scenarios := []struct {
		name   string
		Key    string
		expErr error
		exp    string
	}{
		{"new Key",
			"D",
			ErrKeyNotExist,
			"",
		},
		{
			"old Key",
			"B",
			nil,
			"2",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actual, err := store.Get(scenario.Key)
			if scenario.expErr != nil {
				if err == nil {
					t.Errorf("expected error %v, but got %v instead\n", scenario.expErr, err)
				}
				if !errors.Is(err, scenario.expErr) {
					t.Errorf("expected error %v, but got %v instead\n", scenario.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error %v, but got %v instead\n", scenario.expErr, err)
			}
			if actual != scenario.exp {
				t.Errorf("expected value (%v), but got (%v) instead", scenario.exp, actual)
			}
		})
	}
}

func TestKVStore_Delete(t *testing.T) {
	scenarios := []struct {
		name   string
		Key    string
		expErr error
	}{
		{"new Key",
			"new Key",
			nil,
		},
		{
			"old Key",
			"old Key",
			nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := store.Delete(scenario.Key)
			if scenario.expErr != nil {
				if err == nil {
					t.Errorf("expected error %v, but got %v instead\n", scenario.expErr, err)
				}
				if !errors.Is(err, scenario.expErr) {
					t.Errorf("expected error %v, but got %v instead\n", scenario.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error %v, but got %v instead\n", scenario.expErr, err)
			}
		})
	}
}

func setup() error {
	store = NewKVStore()
	store.s["A"] = "1"
	store.s["B"] = "2"
	store.s["C"] = "3"
	return nil
}

func tearDown() {
	store = nil
}
