package kvstore

import (
	"errors"
	"sync"
)

var ErrNoSuchKey = errors.New("no such key")

type Storage struct {
	sync.RWMutex
	m map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		RWMutex: sync.RWMutex{},
		m:       map[string]string{},
	}
}

func (s *Storage) Put(key, value string) error {
	s.Lock()
	defer s.Unlock()

	s.m[key] = value

	return nil
}

func (s *Storage) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.m[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return v, nil
}

func (s *Storage) Delete(key string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)

	return nil
}
