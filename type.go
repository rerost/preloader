package preloader

import "context"

type NodeFunc[NodeID comparable, Node Resource[NodeID]] func(context.Context, []NodeID) ([]Node, error)
type RelationFunc[ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] func(context.Context, []Parent) (map[ParentID][]NodeID, error)
