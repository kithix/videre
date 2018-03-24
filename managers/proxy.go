package managers

import (
	"io"
	"net/http"
	"sync"
)

// Openers should block until the source is ready.
type Opener func() io.Closer

type ProxyManager struct {
	closer io.Closer
	opener Opener

	handler http.Handler

	connections map[*http.Request]struct{}
	connLock    sync.Mutex
}

func (m *ProxyManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.startConnection(r)
	defer m.endConnection(r)
	m.handler.ServeHTTP(w, r)
}

func (m *ProxyManager) endConnection(uniqueID *http.Request) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	delete(m.connections, uniqueID)
	if len(m.connections) == 0 {
		m.closer.Close()
	}
}

// Are all requests pointers unique? probably for their lifetime.
func (m *ProxyManager) startConnection(uniqueID *http.Request) {
	m.connLock.Lock()
	defer m.connLock.Unlock()
	if len(m.connections) == 0 {
		m.closer = m.opener()
	}

	m.connections[uniqueID] = struct{}{}
}

func NewProxyManager(opener Opener, handler http.Handler) *ProxyManager {
	return &ProxyManager{
		opener:      opener,
		handler:     handler,
		connections: make(map[*http.Request]struct{}),
	}
}
