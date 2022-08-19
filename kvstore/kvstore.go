package kvstore

import "errors"

var ErrNoSuchKey = errors.New("no such key")

type Storage struct {
	m map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		m: map[string]string{},
	}
}

func (s *Storage) Put(key, value string) error {
	s.m[key] = value

	return nil
}

func (s *Storage) Get(key string) (string, error) {
	v, ok := s.m[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return v, nil
}

func (s *Storage) Delete(key string) error {
	delete(s.m, key)

	return nil
}
