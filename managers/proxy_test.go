package managers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockCloser struct {
	b bool
	c chan struct{}
}

func (m *MockCloser) Close() error {
	m.b = false
	return nil
}

func (m *MockCloser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.b == true {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
	m.c <- struct{}{}
}

func NewMockCloserMaker() (Opener, *MockCloser) {
	mockCloser := &MockCloser{
		b: false,
		c: make(chan struct{}),
	}
	return func() io.Closer {
		mockCloser.b = true
		return mockCloser
	}, mockCloser
}

func testProxyManagerServeHTTP(s *ProxyManager) bool {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://target.url", nil)

	go s.ServeHTTP(w, r)
	time.Sleep(1 * time.Microsecond)

	if w.Code != 200 {
		return false
	}
	return true
}
func closeHTTP(mockCloser *MockCloser) {
	<-mockCloser.c
	time.Sleep(1 * time.Microsecond)
}

func TestProxyManager(t *testing.T) {
	closerMaker, mockCloser := NewMockCloserMaker()
	s := NewProxyManager(closerMaker, mockCloser)

	testProxyManagerServeHTTP(s)
	if !mockCloser.b {
		t.Error("Did not correctly call closerMaker")
		t.Fail()
	}
	closeHTTP(mockCloser)
	if mockCloser.b {
		t.Error("Did not correctly call close on mockCloser")
		t.Fail()
	}
}

func TestProxyManagerMultipleConnections(t *testing.T) {
	closerMaker, mockCloser := NewMockCloserMaker()
	s := NewProxyManager(closerMaker, mockCloser)

	testProxyManagerServeHTTP(s)
	testProxyManagerServeHTTP(s)
	closeHTTP(mockCloser)
	if !mockCloser.b {
		t.Error("Incorrectly called close on mockCloser")
		t.Fail()
	}
	closeHTTP(mockCloser)
	if mockCloser.b {
		t.Error("Did not correctly call close on mockCloser")
		t.Fail()
	}
}
