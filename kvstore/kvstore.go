package kvstore

import (
	"errors"
	"sync"
)

var ErrNoSuchKey = errors.New("no such key")

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)

	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type Storage struct {
	sync.RWMutex
	m map[string]string
}

var store = Storage{
	m:       map[string]string{},
	RWMutex: sync.RWMutex{},
}

// func NewStorage() *Storage {
// 	return &Storage{
// 		RWMutex: sync.RWMutex{},
// 		m:       map[string]string{},
// 	}
// }

func Put(key, value string) error {
	store.Lock()
	defer store.Unlock()

	store.m[key] = value

	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()

	v, ok := store.m[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return v, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)

	return nil
}
