package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	Mutex        sync.Mutex
	ReverseProxy *httputil.ReverseProxy
}
func (b *Backend) SetAlive() {
	b.Mutex.Lock()

	b.Alive = true

	b.Mutex.Unlock()
}
func (b *Backend) isAlive() bool {
	defer b.Mutex.Lock()

	b.Mutex.Unlock()

	isAlive := b.Alive

	return isAlive
}

type ServerPool struct {
	Backend []*Backend
	current uint64
}
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.Backend)))
}
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.Backend) + next

	for i := next; i < l; i++ {
		idx := i % len(s.Backend)

		if s.Backend[idx].Alive {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}

			return s.Backend[idx]
		}
	}

	return nil
}

func main() {
	// url, _ := url.Parse("http://localhost:8080")

	// reverseProxy := httputil.NewSingleHostReverseProxy(url)

}
