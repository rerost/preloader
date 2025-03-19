package preloader

import (
	"context"
)

// LoadableKey is a type used to create unique keys for loadable types
type LoadableKey string

// GetLoadableFromProvider is a helper function to get a Loadable from a provider
func GetLoadableFromProvider[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	provider LoadableProvider,
	key string,
) Loadable[ParentID, Parent, NodeID, Node] {
	loadable := provider.GetLoadable(key)
	if loadable == nil {
		return EmptyLoadable[ParentID, Parent, NodeID, Node]()
	}
	
	typedLoadable, ok := loadable.(Loadable[ParentID, Parent, NodeID, Node])
	if !ok {
		return EmptyLoadable[ParentID, Parent, NodeID, Node]()
	}
	
	return typedLoadable
}

// GetHasOneLoadableFromProvider is a helper function to get a HasOneLoadable from a provider
func GetHasOneLoadableFromProvider[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	provider LoadableProvider,
	key string,
) HasOneLoadable[ParentID, Parent, NodeID, Node] {
	loadable := provider.GetLoadable(key)
	if loadable == nil {
		return EmptyHasOneLoadable[ParentID, Parent, NodeID, Node]()
	}
	
	typedLoadable, ok := loadable.(HasOneLoadable[ParentID, Parent, NodeID, Node])
	if !ok {
		return EmptyHasOneLoadable[ParentID, Parent, NodeID, Node]()
	}
	
	return typedLoadable
}
