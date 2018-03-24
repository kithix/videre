package stores

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

var (
	ErrorNothingToFetch = errors.New("No stored bytes to fetch")
)

type SingleStore struct {
	lock        sync.Mutex
	storedBytes []byte
}

func (s *SingleStore) Write(p []byte) (int, error) {
	s.lock.Lock()

	s.storedBytes = make([]byte, len(p))
	i := copy(s.storedBytes, p)

	s.lock.Unlock()
	return i, nil
}

func (s *SingleStore) Fetch() (io.Reader, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.storedBytes == nil {
		return nil, ErrorNothingToFetch
	}
	newBytes := make([]byte, len(s.storedBytes))
	copy(newBytes, s.storedBytes)
	return bytes.NewReader(newBytes), nil
}

func NewSingleStore() *SingleStore {
	return &SingleStore{}
}
