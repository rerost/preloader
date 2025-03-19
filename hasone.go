package preloader

import (
	"context"
	"fmt"

	"github.com/rerost/preloader/util"
)

type HasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] interface {
	Load(ctx context.Context, parent Parent) (Node, error)
	Child(child ...ChildLoadable[Node]) HasOneLoadable[ParentID, Parent, NodeID, Node]
	Preload(ctx context.Context, parents []Parent) error

	GetLoaded() ([]Node, error)
}

func NewHasOneLoadable[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]](
	typeName string,
	idsLoader RelationFunc[ParentID, Parent, NodeID, Node],
	loader NodeFunc[NodeID, Node],
	acceptNotFound bool,
) HasOneLoadable[ParentID, Parent, NodeID, Node] {
	return &hasOneLoadableImpl[ParentID, Parent, NodeID, Node]{
		typeName:       typeName,
		loaded:         false,
		values:         util.SyncMap[ParentID, Node]{},
		idLoader:       idsLoader,
		loader:         loader,
		acceptNotFound: acceptNotFound,
	}
}

type hasOneLoadableImpl[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct {
	loaded         bool
	values         util.SyncMap[ParentID, Node]
	idLoader       RelationFunc[ParentID, Parent, NodeID, Node]
	loader         NodeFunc[NodeID, Node]
	children       []ChildLoadable[Node]
	typeName       string
	acceptNotFound bool
}

func (l *hasOneLoadableImpl[ParentID, Parent, NodeID, Node]) Load(ctx context.Context, parent Parent) (Node, error) {
	loadedResource, ok := l.values.Load(parent.GetResourceID())
	if ok {
		return *loadedResource, nil
	}

	var empty Node
	ids, err := l.idLoader(ctx, []Parent{parent})
	if err != nil {
		return empty, err
	}

	res, err := l.loader(ctx, ids[parent.GetResourceID()])
	if err != nil {
		return empty, err
	}

	l.values.Store(parent.GetResourceID(), res[0])
	return res[0], nil
}

func (l *hasOneLoadableImpl[ParentID, Parent, NodeID, Node]) Child(child ...ChildLoadable[Node]) HasOneLoadable[ParentID, Parent, NodeID, Node] {
	l.children = append(l.children, child...)
	return l
}

func (l *hasOneLoadableImpl[ParentID, Parent, NodeID, Node]) Preload(ctx context.Context, parents []Parent) error {
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

		if len(parentsNode) > 1 {
			return TooManyHasOneError
		}
		if len(parentsNode) == 0 {
			if l.acceptNotFound {
				continue
			}
			return NotFoundHasOneError
		}
		if len(parentsNode) == 1 {
			l.values.Store(currentParent.GetResourceID(), parentsNode[0])
		}
	}

	return nil
}

var TooManyHasOneError = fmt.Errorf("TooManyHasOne")
var NotFoundHasOneError = fmt.Errorf("NotFoundHasOne")

func (l *hasOneLoadableImpl[ParentID, Parent, NodeID, Node]) GetLoaded() ([]Node, error) {
	res := []Node{}
	l.values.Range(func(_ ParentID, v Node) bool {
		res = append(res, v)
		return true
	})

	return res, nil
}
