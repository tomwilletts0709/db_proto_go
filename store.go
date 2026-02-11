package main

import (
	"sync"
)

type Store struct {
	sets  map[string]string
	hsets map[string]map[string]string
	mu    sync.RWMutex
	aof   *Aof
}

func NewStore(aof *Aof) *Store {
	return &Store{
		sets:  make(map[string]string),
		hsets: make(map[string]map[string]string),
		mu:    sync.RWMutex{},
		aof:   aof}
}

func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	s.sets[key] = value
	s.mu.Unlock()
}

func (s *Store) Get(key string) string {
	s.mu.RLock()
	value, ok := s.sets[key]
	s.mu.RUnlock()
	return value, ok
}

func (s *Store) Hset(hash string, key string, value string) {
	s.mu.Lock()
	if _, ok := s.hsets[hash]; !ok {
		s.hsets[hash] = map[string]string{}
	}
	s.hsets[hash][key] = value
	s.mu.Unlock()
}
