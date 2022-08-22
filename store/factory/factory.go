package factory

import (
	"fmt"
	"sync"

	"bookstore/store"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]store.Store)
)

func Register(name string, p store.Store) {
	providersMu.Lock()
	defer providersMu.Unlock()

	if p == nil {
		panic("store: Register provider is nil")
	}

	if _, ok := providers[name]; ok {
		panic("store: Register called twice for provider: " + name)
	}
	providers[name] = p
}

func New(name string) (store.Store, error) {
	providersMu.RLock()
	p, ok := providers[name]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown provider %s", name)
	}
	return p, nil
}
