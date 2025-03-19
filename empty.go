package preloader

import (
	"context"
	"fmt"
)

func EmptyHasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]]() HasOneLoadable[ParentID, Parent, NodeID, Node] {
	return &emptyHasOneLoadableImpl[ParentID, Parent, NodeID, Node]{}
}

type emptyHasOneLoadableImpl[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct{}

func (l *emptyHasOneLoadableImpl[ParentID, Parent, NodeID, Node]) Load(ctx context.Context, parent Parent) (Node, error) {
	var empty Node
	return empty, NoLoadableError
}

func (l *emptyHasOneLoadableImpl[ParentID, Parent, NodeID, Node]) Child(child ...ChildLoadable[Node]) HasOneLoadable[ParentID, Parent, NodeID, Node] {
	return l
}

func (l *emptyHasOneLoadableImpl[ParentID, Parent, NodeID, Node]) Preload(ctx context.Context, parents []Parent) error {
	return nil
}

func (l *emptyHasOneLoadableImpl[ParentID, Parent, NodeID, Node]) GetLoaded() ([]Node, error) {
	return nil, nil
}

func EmptyLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]]() Loadable[ParentID, Parent, NodeID, Node] {
	return &emptyLoadableImpl[ParentID, Parent, NodeID, Node]{}
}

type emptyLoadableImpl[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct {
	loadableImpl[ParentID, Parent, NodeID, Node]
}

func (l *emptyLoadableImpl[ParentID, Parent, NodeID, Node]) Load(ctx context.Context, parent Parent) ([]Node, error) {
	return nil, nil
}

func (l *emptyLoadableImpl[ParentID, Parent, NodeID, Node]) Child(child ...ChildLoadable[Node]) Loadable[ParentID, Parent, NodeID, Node] {
	return l
}

func (l *emptyLoadableImpl[ParentID, Parent, NodeID, Node]) Preload(ctx context.Context, parents []Parent) error {
	return nil
}

var NoLoadableError = fmt.Errorf("NoLoadable")
