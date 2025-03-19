package preloader

import (
	"context"

	"github.com/rerost/preloader/util"
)

func NewLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	typeName string,
	idsLoader RelationFunc[ParentID, Parent, NodeID, Node],
	loader NodeFunc[NodeID, Node],
) Loadable[ParentID, Parent, NodeID, Node] {
	return &loadableImpl[ParentID, Parent, NodeID, Node]{
		typeName: typeName,
		loaded:   false,
		values:   util.SyncMap[ParentID, []Node]{},
		idLoader: idsLoader,
		loader:   loader,
	}
}

type Resource[TID comparable] interface {
	GetResourceID() TID
}

type ChildLoadable[Parent any] interface {
	Preload(ctx context.Context, parent []Parent) error
}

type Loadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] interface {
	Load(ctx context.Context, parent Parent) ([]Node, error)
	Child(child ...ChildLoadable[Node]) Loadable[ParentID, Parent, NodeID, Node]
	Preload(ctx context.Context, parents []Parent) error

	GetLoaded() ([]Node, error)
}

type loadableImpl[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct {
	loaded   bool
	values   util.SyncMap[ParentID, []Node]
	idLoader RelationFunc[ParentID, Parent, NodeID, Node]
	loader   NodeFunc[NodeID, Node]
	children []ChildLoadable[Node]
	typeName string
}

func (l *loadableImpl[ParentID, Parent, NodeID, Node]) Load(ctx context.Context, parent Parent) ([]Node, error) {
	loadedResource, ok := l.values.Load(parent.GetResourceID())
	if ok {
		return *loadedResource, nil
	}

	tids, err := l.idLoader(ctx, []Parent{parent})
	if err != nil {
		return nil, err
	}

	res, err := l.loader(ctx, tids[parent.GetResourceID()])
	if err != nil {
		return nil, err
	}

	l.values.Store(parent.GetResourceID(), res)
	return res, nil
}

func (l *loadableImpl[ParentID, Parent, NodeID, Node]) Child(child ...ChildLoadable[Node]) Loadable[ParentID, Parent, NodeID, Node] {
	l.children = append(l.children, child...)
	return l
}

func (l *loadableImpl[ParentID, Parent, NodeID, Node]) Preload(ctx context.Context, parents []Parent) error {
	idsLoader := l.idLoader
	loader := l.loader

	// Parent -> []NodeID
	relations, err := idsLoader(ctx, parents)
	if err != nil {
		return err
	}

	// 参照されているNodeの一覧。親子関係は無視している
	var nodes []Node
	{
		flatten := []NodeID{}
		for _, nodeID := range relations {
			flatten = append(flatten, nodeID...)
		}

		nodes, err = loader(ctx, flatten)
		if err != nil {
			return err
		}
	}

	// ParentID -> Parent
	parentIDMap := make(map[ParentID]Parent, len(parents))
	for _, parent := range parents {
		parentIDMap[parent.GetResourceID()] = parent
	}

	// NodeID -> Node
	nodeIDsMap := make(map[NodeID]Node, len(nodes))
	for _, node := range nodes {
		nodeIDsMap[node.GetResourceID()] = node
	}

	for _, currentParent := range parents {
		nodeIDs := relations[currentParent.GetResourceID()]
		// currentParent -> Node
		parentsNode := make([]Node, 0, len(nodeIDs))
		for _, nodeID := range nodeIDs {
			parentsNode = append(parentsNode, nodeIDsMap[nodeID])
		}

		l.values.Store(currentParent.GetResourceID(), parentsNode)
	}

	for _, child := range l.children {
		err := child.Preload(ctx, nodes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *loadableImpl[ParentID, Parent, NodeID, Node]) GetLoaded() ([]Node, error) {
	res := []Node{}
	l.values.Range(func(_ ParentID, vs []Node) bool {
		res = append(res, vs...)
		return true
	})

	return res, nil
}

func Preload[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	ctx context.Context,
	parents []Parent,
	loadables ...Loadable[ParentID, Parent, NodeID, Node],
) error {
	if len(parents) == 0 {
		return nil
	}
	if len(loadables) == 0 {
		return nil
	}

	for _, loadable := range loadables {
		loadable.Preload(ctx, parents)
	}

	return nil
}
