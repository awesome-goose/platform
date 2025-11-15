// (c) https://github.com/golobby/container/blob/v3/pkg/container/container.go

package core

import (
	"errors"
	"fmt"
	"reflect"
)

type binding struct {
	resolver any
	instance any
}

func NewBinding(resolver any, instance any) binding {
	return binding{
		resolver: resolver,
		instance: instance,
	}
}

func (b binding) Resolve(c container) (any, error) {
	if b.instance != nil {
		return b.instance, nil
	}

	return c.invoke(b.resolver)
}

type container map[reflect.Type]map[string]binding

func NewContainer() container {
	return make(container)
}

func (c container) Register(resolver any, name string, singleton bool) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("container: the resolver must be a function")
	}

	for i := 0; i < reflectedResolver.NumOut(); i++ {
		if _, exist := c[reflectedResolver.Out(i)]; !exist {
			c[reflectedResolver.Out(i)] = make(map[string]binding)
		}

		if singleton {
			instance, err := c.invoke(resolver)
			if err != nil {
				return err
			}

			c[reflectedResolver.Out(i)][name] = binding{resolver: resolver, instance: instance}
		} else {
			c[reflectedResolver.Out(i)][name] = binding{resolver: resolver}
		}
	}

	return nil
}

func (c container) Resolve(abstraction any, name string) error {
	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return errors.New("container: invalid abstraction")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()

		if concrete, exist := c[elem][name]; exist {
			if instance, err := concrete.Resolve(c); err == nil {
				reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
				return nil
			} else {
				return err
			}
		}

		return errors.New("container: no concrete found for: " + elem.String())
	}

	return errors.New("container: invalid abstraction")
}

func (c container) Make(value any) (any, error) {
	name := c.getTypeName(value)
	t := reflect.TypeOf(value)

	// Dereference pointer types to get underlying type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Try to find "New" method
	// or the first method that starts with "New"
	constructor := reflect.ValueOf(t).MethodByName("New")
	if !constructor.IsValid() {
		// Accept function directly
		constructor = reflect.ValueOf(t)
		if constructor.Kind() != reflect.Func || !constructor.IsValid() {
			return nil, fmt.Errorf("container: no constructor found for %T", t)
		}
	}

	ctorType := constructor.Type()
	args := make([]reflect.Value, ctorType.NumIn())

	for i := 0; i < ctorType.NumIn(); i++ {
		argType := ctorType.In(i)

		if c.isPrimitive(argType) {
			args[i] = reflect.Zero(argType)
			continue
		}

		// Try to resolve from container or recursively Make
		if b, ok := c[argType][name]; ok {
			resolved, err := b.Resolve(c)
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(resolved)
		} else {
			// Use Make recursively
			val := reflect.New(argType).Interface()
			resolved, err := c.Make(val)
			if err != nil {
				return nil, fmt.Errorf("container: cannot resolve %v: %w", argType, err)
			}
			args[i] = reflect.ValueOf(resolved)
		}
	}

	results := constructor.Call(args)
	if len(results) == 0 {
		return nil, fmt.Errorf("container: constructor did not return anything")
	}

	var instance reflect.Value
	if len(results) == 2 {
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}
		instance = results[0]
	} else {
		instance = results[0]
	}

	// Ensure it's a pointer
	instanceType := instance.Type()
	if instanceType.Kind() != reflect.Ptr {
		ptr := reflect.New(instanceType)
		ptr.Elem().Set(instance)
		instance = ptr
	}

	// Register as singleton
	finalType := instance.Type()
	if _, exists := c[finalType]; !exists {
		c[finalType] = make(map[string]binding)
	}
	c[finalType][name] = binding{resolver: nil, instance: instance.Interface()}

	return instance.Interface(), nil
}

func (c container) invoke(function any) (any, error) {
	args, err := c.arguments(function)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(function).NumOut() == 1 {
		return reflect.ValueOf(function).Call(args)[0].Interface(), nil
	} else if reflect.TypeOf(function).NumOut() == 2 {
		values := reflect.ValueOf(function).Call(args)
		return values[0].Interface(), values[1].Interface().(error)
	}

	return nil, errors.New("container: resolver function signature is invalid")
}

func (c container) arguments(function any) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)

		if concrete, exist := c[abstraction][""]; exist {
			instance, _ := concrete.Resolve(c)

			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, errors.New("container: no concrete found for: " + abstraction.String())
		}
	}

	return arguments, nil
}

func (c container) isPrimitive(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func (c container) getTypeName(v any) string {
	return reflect.TypeOf(v).String()
}
