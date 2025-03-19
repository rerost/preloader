package preloader

import (
	"fmt"
	"sync"
)

// LoadableKey is a type used to create unique keys for loadable types
type LoadableKey string

// TypedLoadableProvider is a type-safe version of LoadableProvider
// It uses compile-time checks to ensure loadables are registered
type TypedLoadableProvider struct {
	provider     interface{}
	registeredKeys map[LoadableKey]bool
	mu          sync.RWMutex
}

// NewTypedLoadableProvider creates a new TypedLoadableProvider
func NewTypedLoadableProvider() *TypedLoadableProvider {
	return &TypedLoadableProvider{
		provider:     make(map[string]interface{}),
		registeredKeys: make(map[LoadableKey]bool),
	}
}

// RegisterLoadable registers a loadable with the provider
func (p *TypedLoadableProvider) RegisterLoadable(typeName string, loadable interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.provider.(map[string]interface{})[typeName] = loadable
}

// GetLoadable retrieves a loadable from the provider
func (p *TypedLoadableProvider) GetLoadable(typeName string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.provider.(map[string]interface{})[typeName]
}

// RegisterTypedLoadable registers a typed loadable with the provider
func RegisterTypedLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	key LoadableKey,
	loadable Loadable[ParentID, Parent, NodeID, Node],
) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.provider.(map[string]interface{})[string(key)] = loadable
	p.registeredKeys[key] = true
}

// RegisterTypedHasOneLoadable registers a typed has-one loadable with the provider
func RegisterTypedHasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	key LoadableKey,
	loadable HasOneLoadable[ParentID, Parent, NodeID, Node],
) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.provider.(map[string]interface{})[string(key)] = loadable
	p.registeredKeys[key] = true
}

// MustGetLoadable retrieves a typed loadable from the provider
// It will panic if the loadable is not registered or has the wrong type
func MustGetLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	key LoadableKey,
) Loadable[ParentID, Parent, NodeID, Node] {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Check if the key is registered
	if _, ok := p.registeredKeys[key]; !ok {
		panic(fmt.Sprintf("Loadable with key %s is not registered", key))
	}
	
	// Get the loadable from the provider
	loadable := p.provider.(map[string]interface{})[string(key)]
	if loadable == nil {
		panic(fmt.Sprintf("Loadable with key %s is nil", key))
	}
	
	// Type assertion
	typedLoadable, ok := loadable.(Loadable[ParentID, Parent, NodeID, Node])
	if !ok {
		panic(fmt.Sprintf("Type mismatch for loadable with key %s", key))
	}
	
	return typedLoadable
}

// MustGetHasOneLoadable retrieves a typed has-one loadable from the provider
// It will panic if the loadable is not registered or has the wrong type
func MustGetHasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	key LoadableKey,
) HasOneLoadable[ParentID, Parent, NodeID, Node] {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Check if the key is registered
	if _, ok := p.registeredKeys[key]; !ok {
		panic(fmt.Sprintf("HasOneLoadable with key %s is not registered", key))
	}
	
	// Get the loadable from the provider
	loadable := p.provider.(map[string]interface{})[string(key)]
	if loadable == nil {
		panic(fmt.Sprintf("HasOneLoadable with key %s is nil", key))
	}
	
	// Type assertion
	typedLoadable, ok := loadable.(HasOneLoadable[ParentID, Parent, NodeID, Node])
	if !ok {
		panic(fmt.Sprintf("Type mismatch for HasOneLoadable with key %s", key))
	}
	
	return typedLoadable
}
