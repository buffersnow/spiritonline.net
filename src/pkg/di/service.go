package di

import (
	"errors"
	"fmt"
	"reflect"
)

func (c *Container) Register(service any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := reflect.ValueOf(service)
	typ := val.Type()

	// Check if it's a zero-arg constructor: func() (*T, error)
	if typ.Kind() == reflect.Func &&
		typ.NumIn() == 0 &&
		typ.NumOut() == 2 &&
		containsError(typ) {

		results := val.Call(nil)

		if !results[1].IsNil() {
			err := results[1].Interface().(error)
			return fmt.Errorf("di: %w", err)
		}

		instance := results[0]
		instanceType := instance.Type()

		if _, exists := c.instances[instanceType]; exists {
			return errors.New("di: service already registered")
		}

		c.instances[instanceType] = instance
		return nil
	}

	// Standard: registering a concrete instance directly
	if _, exists := c.instances[typ]; exists {
		return errors.New("di: service already registered")
	}
	c.instances[typ] = val
	return nil
}

func (c *Container) Locate(target any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ptr := reflect.ValueOf(target)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		return errors.New("di: locatable target must be a non-nil pointer")
	}

	targetType := ptr.Elem().Type()

	if instance, ok := c.instances[targetType]; ok {
		ptr.Elem().Set(instance)
		return nil
	}

	if provider, ok := c.providers[targetType]; ok {
		instance, err := c.invokeProvider(provider)
		if err != nil {
			return fmt.Errorf("di: %w", err)
		}
		ptr.Elem().Set(instance)
		// Cache the instance
		c.instances[targetType] = instance
		return nil
	}

	return errors.New("no instance or provider found")
}

// invokeProvider executes a provider function and checks for returned error
func (c *Container) invokeProvider(provider reflect.Value) (reflect.Value, error) {
	args, err := c.invokeWithDeps(provider)
	if err != nil {
		return reflect.Value{}, err
	}
	results := provider.Call(args)

	if !results[1].IsNil() {
		return reflect.Value{}, results[1].Interface().(error)
	}
	return results[0], nil
}
