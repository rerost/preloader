package dynoproxy

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ovechkin-dm/go-dyno/pkg/dyno"
	"github.com/rerost/preloader"
)

func DynamicLoadable[ParentID comparable, Parent preloader.Resource[ParentID], NodeID comparable, Node preloader.Resource[NodeID]](
	typeName string,
	idsLoader preloader.RelationFunc[ParentID, Parent, NodeID, Node],
	loader preloader.NodeFunc[NodeID, Node],
) (preloader.Loadable[ParentID, Parent, NodeID, Node], error) {
	loadable := preloader.NewLoadable(typeName, idsLoader, loader)
	
	handler := &LoadableHandler[ParentID, Parent, NodeID, Node]{
		Impl: loadable,
	}
	
	return dyno.Dynamic[preloader.Loadable[ParentID, Parent, NodeID, Node]](handler.Handle)
}

func DynamicHasOneLoadable[ParentID comparable, Parent preloader.Resource[ParentID], NodeID comparable, Node preloader.Resource[NodeID]](
	typeName string,
	idsLoader preloader.RelationFunc[ParentID, Parent, NodeID, Node],
	loader preloader.NodeFunc[NodeID, Node],
	acceptNotFound bool,
) (preloader.HasOneLoadable[ParentID, Parent, NodeID, Node], error) {
	loadable := preloader.NewHasOneLoadable(typeName, idsLoader, loader, acceptNotFound)
	
	handler := &HasOneLoadableHandler[ParentID, Parent, NodeID, Node]{
		Impl: loadable,
	}
	
	return dyno.Dynamic[preloader.HasOneLoadable[ParentID, Parent, NodeID, Node]](handler.Handle)
}

type LoadableHandler[ParentID comparable, Parent preloader.Resource[ParentID], NodeID comparable, Node preloader.Resource[NodeID]] struct {
	Impl preloader.Loadable[ParentID, Parent, NodeID, Node]
}

func (h *LoadableHandler[ParentID, Parent, NodeID, Node]) Handle(m reflect.Method, values []reflect.Value) []reflect.Value {
	return reflect.ValueOf(h.Impl).MethodByName(m.Name).Call(values)
}

type HasOneLoadableHandler[ParentID comparable, Parent preloader.Resource[ParentID], NodeID comparable, Node preloader.Resource[NodeID]] struct {
	Impl preloader.HasOneLoadable[ParentID, Parent, NodeID, Node]
}

func (h *HasOneLoadableHandler[ParentID, Parent, NodeID, Node]) Handle(m reflect.Method, values []reflect.Value) []reflect.Value {
	return reflect.ValueOf(h.Impl).MethodByName(m.Name).Call(values)
}

func AutoInject(target interface{}, loadables map[string]interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	
	targetElem := targetValue.Elem()
	if targetElem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	
	targetType := targetElem.Type()
	for i := 0; i < targetElem.NumField(); i++ {
		field := targetElem.Field(i)
		fieldType := targetType.Field(i)
		
		if !field.CanSet() {
			continue
		}
		
		fieldName := fieldType.Name
		loadable, ok := loadables[fieldName]
		if !ok {
			continue
		}
		
		loadableValue := reflect.ValueOf(loadable)
		if !loadableValue.Type().AssignableTo(field.Type()) {
			return fmt.Errorf("loadable type %s is not assignable to field %s", loadableValue.Type(), fieldName)
		}
		
		field.Set(loadableValue)
	}
	
	return nil
}

type ModelBuilder struct {
	loadables map[string]interface{}
}

func NewModelBuilder() *ModelBuilder {
	return &ModelBuilder{
		loadables: make(map[string]interface{}),
	}
}

func (b *ModelBuilder) AddLoadable(name string, loadable interface{}) *ModelBuilder {
	b.loadables[name] = loadable
	return b
}

func (b *ModelBuilder) Build(target interface{}) error {
	return AutoInject(target, b.loadables)
}

type PreloadHelper struct {
	ctx      context.Context
	loadables []interface{}
}

func NewPreloadHelper(ctx context.Context) *PreloadHelper {
	return &PreloadHelper{
		ctx:      ctx,
		loadables: make([]interface{}, 0),
	}
}

func (h *PreloadHelper) AddLoadable(loadable interface{}) *PreloadHelper {
	h.loadables = append(h.loadables, loadable)
	return h
}

func (h *PreloadHelper) Preload(parents interface{}) error {
	parentsValue := reflect.ValueOf(parents)
	if parentsValue.Kind() != reflect.Slice {
		return fmt.Errorf("parents must be a slice")
	}
	
	for _, loadable := range h.loadables {
		loadableValue := reflect.ValueOf(loadable)
		preloadMethod := loadableValue.MethodByName("Preload")
		if !preloadMethod.IsValid() {
			return fmt.Errorf("loadable does not have Preload method")
		}
		
		args := []reflect.Value{reflect.ValueOf(h.ctx), parentsValue}
		results := preloadMethod.Call(args)
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().(error)
		}
	}
	
	return nil
}
