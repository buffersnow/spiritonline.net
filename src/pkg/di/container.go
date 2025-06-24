package di

import (
	"errors"
	"fmt"
	"reflect"
)

func (c *Container) Provide(provider any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := reflect.ValueOf(provider)
	typ := val.Type()

	if typ.Kind() != reflect.Func || typ.NumOut() != 2 {
		return errors.New("di: provider must be a function that has two return values")
	}

	if !containsError(typ) {
		return errors.New("di: second return value of provider must be error")
	}

	outType := typ.Out(0)
	if _, exists := c.providers[outType]; exists {
		return errors.New("di: provider is already registered")
	}
	c.providers[outType] = val
	return nil
}

func (c *Container) Invoke(fn any) error {
	val := reflect.ValueOf(fn)

	args, err := c.invokeWithDeps(val)
	if err != nil {
		return fmt.Errorf("di: %w", err)
	}
	val.Call(args)
	return nil
}

func (c *Container) invokeWithDeps(fn reflect.Value) ([]reflect.Value, error) {
	typ := fn.Type()
	args := make([]reflect.Value, typ.NumIn())

	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)

		c.mu.RLock()
		instance, ok := c.instances[argType]
		provider, hasProvider := c.providers[argType]
		c.mu.RUnlock()

		if !ok {
			if !hasProvider {
				return nil, fmt.Errorf("missing dependency of type %v", argType)
			}
			// Recursively resolve the provider
			instance, err := c.invokeProvider(provider)
			if err != nil {
				return nil, fmt.Errorf("failed to invoke provider for %v: %w", argType, err)
			}
			c.mu.Lock()
			c.instances[argType] = instance
			c.mu.Unlock()
		}

		args[i] = instance
	}
	return args, nil
}
