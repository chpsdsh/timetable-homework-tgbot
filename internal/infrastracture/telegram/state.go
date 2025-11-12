package telegram

import "sync"

type StateStore interface {
	Get(chatID int64) string
	Set(chatID int64, state string)
	Del(chatID int64)
}

type memoryState struct {
	mu sync.RWMutex
	m  map[int64]string
}

func NewMemoryState() StateStore { return &memoryState{m: map[int64]string{}} }

func (s *memoryState) Get(id int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[id]
}
func (s *memoryState) Set(id int64, st string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[id] = st
}
func (s *memoryState) Del(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, id)
}

func NewMemState() StateStore { return &memoryState{m: map[int64]string{}} }
