package preloader

import (
	"sync"
)

// Registered is a phantom type used to mark loadables as registered
type Registered struct{}

// NotRegistered is a phantom type used to mark loadables as not registered
type NotRegistered struct{}

// LoadableKey is a type used to create unique keys for loadable types
type LoadableKey string

// TypedLoadableProvider is a type-safe version of LoadableProvider
// It uses compile-time checks to ensure loadables are registered
type TypedLoadableProvider struct {
	loadables map[string]interface{}
	mu        sync.RWMutex
}

// NewTypedLoadableProvider creates a new TypedLoadableProvider
func NewTypedLoadableProvider() *TypedLoadableProvider {
	return &TypedLoadableProvider{
		loadables: make(map[string]interface{}),
	}
}

// RegisteredLoadable is a type that represents a registered loadable
// The R type parameter is used to encode the registration status
type RegisteredLoadable[R any, ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct {
	Key      LoadableKey
	Loadable Loadable[ParentID, Parent, NodeID, Node]
}

// RegisteredHasOneLoadable is a type that represents a registered has-one loadable
// The R type parameter is used to encode the registration status
type RegisteredHasOneLoadable[R any, ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct {
	Key      LoadableKey
	Loadable HasOneLoadable[ParentID, Parent, NodeID, Node]
}

// RegisterLoadable registers a loadable and returns a RegisteredLoadable with Registered type
func RegisterLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	key LoadableKey,
	loadable Loadable[ParentID, Parent, NodeID, Node],
) RegisteredLoadable[Registered, ParentID, Parent, NodeID, Node] {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.loadables[string(key)] = loadable
	return RegisteredLoadable[Registered, ParentID, Parent, NodeID, Node]{
		Key:      key,
		Loadable: loadable,
	}
}

// RegisterHasOneLoadable registers a has-one loadable and returns a RegisteredHasOneLoadable with Registered type
func RegisterHasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	key LoadableKey,
	loadable HasOneLoadable[ParentID, Parent, NodeID, Node],
) RegisteredHasOneLoadable[Registered, ParentID, Parent, NodeID, Node] {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.loadables[string(key)] = loadable
	return RegisteredHasOneLoadable[Registered, ParentID, Parent, NodeID, Node]{
		Key:      key,
		Loadable: loadable,
	}
}

// GetLoadable retrieves a loadable from the provider
// This function is used internally and should not be called directly
func (p *TypedLoadableProvider) GetLoadable(key string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.loadables[key]
}

// GetRegisteredLoadable retrieves a registered loadable
// This function requires a RegisteredLoadable with Registered type
func GetRegisteredLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	registered RegisteredLoadable[Registered, ParentID, Parent, NodeID, Node],
) Loadable[ParentID, Parent, NodeID, Node] {
	p.mu.RLock()
	defer p.mu.RUnlock()
	loadable := p.loadables[string(registered.Key)]
	return loadable.(Loadable[ParentID, Parent, NodeID, Node])
}

// GetRegisteredHasOneLoadable retrieves a registered has-one loadable
// This function requires a RegisteredHasOneLoadable with Registered type
func GetRegisteredHasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	p *TypedLoadableProvider,
	registered RegisteredHasOneLoadable[Registered, ParentID, Parent, NodeID, Node],
) HasOneLoadable[ParentID, Parent, NodeID, Node] {
	p.mu.RLock()
	defer p.mu.RUnlock()
	loadable := p.loadables[string(registered.Key)]
	return loadable.(HasOneLoadable[ParentID, Parent, NodeID, Node])
}
