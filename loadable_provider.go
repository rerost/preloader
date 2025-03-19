package preloader

import (
	"sync"
)

// LoadableProvider provides a way to associate Loadable objects with models
// without creating circular dependencies between repositories
type LoadableProvider interface {
	// RegisterLoadable registers a Loadable for a specific type
	RegisterLoadable(typeName string, loadable interface{})
	
	// GetLoadable retrieves a Loadable for a specific type
	GetLoadable(typeName string) interface{}
}

// NewLoadableProvider creates a new LoadableProvider
func NewLoadableProvider() LoadableProvider {
	return &loadableProviderImpl{
		loadables: make(map[string]interface{}),
	}
}

type loadableProviderImpl struct {
	mu        sync.RWMutex
	loadables map[string]interface{}
}

func (p *loadableProviderImpl) RegisterLoadable(typeName string, loadable interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.loadables[typeName] = loadable
}

func (p *loadableProviderImpl) GetLoadable(typeName string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.loadables[typeName]
}
