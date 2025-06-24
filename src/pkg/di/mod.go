package di

import (
	"reflect"
	"sync"
)

type Container struct {
	mu        sync.RWMutex
	instances map[reflect.Type]reflect.Value
	providers map[reflect.Type]reflect.Value
	tasks     []func(*Container) error
}

func New() *Container {
	return &Container{
		instances: make(map[reflect.Type]reflect.Value),
		providers: make(map[reflect.Type]reflect.Value),
	}
}

func containsError(typ reflect.Type) bool {
	return typ.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem())
}
