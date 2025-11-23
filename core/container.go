package core

import (
	"reflect"
	"unsafe"

	"github.com/awesome-goose/platform/errors"
)

// ServiceLifecycle defines optional hooks for services
// Services can implement these methods to handle registration and resolution events
// Example:
// func (s *MyService) OnRegister() {}
// func (s *MyService) OnResolve() {}
type ServiceLifecycle interface {
	OnRegister()
	OnResolve()
}

// reflectionCache caches reflect.Type lookups for performance
var reflectionCache = struct {
	byTypeName map[string]reflect.Type
}{
	byTypeName: make(map[string]reflect.Type),
}

type binding struct {
	resolver any
	instance any
}

func (b binding) resolve(c container) (any, error) {
	if b.instance != nil {
		return b.instance, nil
	}

	return c.invoke(b.resolver)
}

type container map[reflect.Type]map[string]binding

// NewContainer creates a new service container.
func NewContainer() container {
	return make(container)
}

func (c container) Register(resolver any, name string, singleton bool) error {
	reflectedResolver := c.getTypeFromCache(name, resolver)
	regName := name
	if regName == "" {
		regName = reflectedResolver.String()
	}
	if reflectedResolver.Kind() == reflect.Func {
		for i := 0; i < reflectedResolver.NumOut(); i++ {
			outType := c.getTypeFromCache(reflectedResolver.Out(i).String(), reflect.Zero(reflectedResolver.Out(i)).Interface())
			// If the return type is an interface, register under the interface type
			if outType.Kind() == reflect.Interface {
				if _, exist := c[outType]; !exist {
					c[outType] = make(map[string]binding)
				}
				if singleton {
					instance, err := c.invoke(resolver)
					if err != nil {
						return err
					}
					c[outType][regName] = binding{resolver: resolver, instance: instance}
				} else {
					c[outType][regName] = binding{resolver: resolver}
				}
				continue
			}
			// Otherwise, register as usual
			if _, exist := c[outType]; !exist {
				c[outType] = make(map[string]binding)
			}
			if singleton {
				instance, err := c.invoke(resolver)
				if err != nil {
					return err
				}
				c[outType][regName] = binding{resolver: resolver, instance: instance}
			} else {
				c[outType][regName] = binding{resolver: resolver}
			}
		}
		return nil
	}
	// Handle struct registration
	if reflectedResolver.Kind() == reflect.Struct {
		var instance any
		var err error
		serviceTypeName := reflectedResolver.Name()
		newMethodNames := []string{"New"}
		if serviceTypeName != "" {
			newMethodNames = append(newMethodNames, "New"+serviceTypeName)
		}
		var foundNew bool
		for _, methodName := range newMethodNames {
			if _, ok := reflectedResolver.MethodByName(methodName); ok {
				foundNew = true
				m := reflect.ValueOf(resolver).MethodByName(methodName)
				results := m.Call(nil)
				if len(results) == 2 {
					if !results[1].IsNil() {
						return results[1].Interface().(error)
					}
					instance = results[0].Interface()
				} else if len(results) == 1 {
					instance = results[0].Interface()
				} else {
					return errors.ErrConstructorDidNotReturnAnything
				}
				break
			}
		}
		if !foundNew {
			instance, err = c.Make(resolver)
			if err != nil {
				return err
			}
		}
		t := c.getTypeFromCache("", instance)
		if _, exist := c[t]; !exist {
			c[t] = make(map[string]binding)
		}
		c[t][regName] = binding{resolver: nil, instance: instance}
		c.callLifecycleHook(instance, "OnRegister")
		return nil
	}
	return errors.ErrInvalidResolver
}

func (c container) Resolve(abstraction any, name string) error {
	receiverType := c.getTypeFromCache(name, abstraction)
	resolveName := name
	if resolveName == "" {
		resolveName = receiverType.String()
	}
	if receiverType == nil {
		return errors.ErrInvalidAbstraction
	}
	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()
		if concrete, exist := c[elem][resolveName]; exist {
			instance, err := concrete.resolve(c)
			if err == nil {
				reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
				c.callLifecycleHook(instance, "OnResolve")
				return nil
			} else {
				return err
			}
		}
		return errors.ErrNoConcreteFound.WithMeta(elem.String())
	}
	return errors.ErrInvalidAbstraction
}

func (c container) Call(target any, methodName string, args ...any) (any, error) {
	v := reflect.ValueOf(target)
	method := v.MethodByName(methodName)
	if !method.IsValid() {
		return nil, errors.ErrInvalidMethod.WithMeta(methodName)
	}
	methodType := method.Type()
	totalArgs := methodType.NumIn()
	callArgs := make([]reflect.Value, totalArgs)
	for i := 0; i < len(args) && i < totalArgs; i++ {
		callArgs[i] = reflect.ValueOf(args[i])
	}
	for i := len(args); i < totalArgs; i++ {
		argType := methodType.In(i)
		argName := argType.String()
		if bindings, ok := c[argType]; ok {
			if binding, ok := bindings[argName]; ok {
				instance, err := binding.resolve(c)
				if err != nil {
					return nil, err
				}
				callArgs[i] = reflect.ValueOf(instance)
				continue
			}
		}
		callArgs[i] = reflect.Zero(argType)
	}
	results := method.Call(callArgs)
	if len(results) == 0 {
		return nil, nil
	}
	if methodType.NumOut() > 1 && results[len(results)-1].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		errVal := results[len(results)-1]
		if !errVal.IsNil() {
			return results[0].Interface(), errVal.Interface().(error)
		}
		return results[0].Interface(), nil
	}
	return results[0].Interface(), nil
}

// Make resolves the given type recursively.
// If a New method exists, it is used as a factory. Otherwise, all properties are auto-resolved (even unexported).
func (c container) Make(value any) (any, error) {
	t := c.getTypeFromCache("", value)
	name := c.getTypeName(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	regName := name
	if regName == "" {
		regName = t.String()
	}
	// Try to find "New" or "New"+ServiceName method
	serviceType := reflect.TypeOf(value)
	serviceTypeName := serviceType.Name()
	newMethodNames := []string{"New"}
	if serviceTypeName != "" {
		newMethodNames = append(newMethodNames, "New"+serviceTypeName)
	}
	var constructor reflect.Value
	var foundNew bool
	for _, methodName := range newMethodNames {
		m := reflect.ValueOf(value).MethodByName(methodName)
		if m.IsValid() {
			constructor = m
			foundNew = true
			break
		}
	}
	if foundNew {
		ctorType := constructor.Type()
		args := make([]reflect.Value, ctorType.NumIn())
		for i := 0; i < ctorType.NumIn(); i++ {
			argType := ctorType.In(i)
			argName := argType.String()
			argTypeCached := c.getTypeFromCache(argName, reflect.Zero(argType).Interface())
			if c.isPrimitive(argTypeCached) {
				args[i] = reflect.Zero(argTypeCached)
				continue
			}
			if b, ok := c[argTypeCached][argName]; ok {
				resolved, err := b.resolve(c)
				if err != nil {
					return nil, err
				}
				args[i] = reflect.ValueOf(resolved)
			} else {
				val := reflect.New(argTypeCached).Interface()
				resolved, err := c.Make(val)
				if err != nil {
					return nil, errors.ErrCannotResolve.WithMeta(argTypeCached).WithError(err)
				}
				args[i] = reflect.ValueOf(resolved)
			}
		}
		results := constructor.Call(args)
		if len(results) == 0 {
			return nil, errors.ErrConstructorDidNotReturnAnything
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
		instanceType := c.getTypeFromCache("", instance.Interface())
		if instanceType.Kind() != reflect.Ptr {
			ptr := reflect.New(instanceType)
			ptr.Elem().Set(instance)
			instance = ptr
		}
		finalType := instance.Type()
		if _, exists := c[finalType]; !exists {
			c[finalType] = make(map[string]binding)
		}
		c[finalType][regName] = binding{resolver: nil, instance: instance.Interface()}
		return instance.Interface(), nil
	}
	if t.Kind() != reflect.Struct {
		return reflect.Zero(t).Interface(), nil
	}
	instance := reflect.New(t)
	v := instance.Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		argTypeCached := c.getTypeFromCache(fieldType.String(), reflect.Zero(fieldType).Interface())
		if c.isPrimitive(argTypeCached) {
			continue
		}
		if argTypeCached.Kind() == reflect.Ptr || argTypeCached.Kind() == reflect.Struct {
			fieldInstance, err := c.Make(reflect.New(argTypeCached).Interface())
			if err != nil {
				return nil, errors.ErrCannotResolve.WithMeta(argTypeCached).WithError(err)
			}
			fv := v.Field(i)
			if fv.CanSet() {
				if argTypeCached.Kind() == reflect.Ptr {
					fv.Set(reflect.ValueOf(fieldInstance))
				} else {
					fv.Set(reflect.ValueOf(fieldInstance).Elem())
				}
			} else {
				// WARNING: Using unsafe to set unexported fields. This may break in future Go versions.
				ptrToField := reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem()
				if argTypeCached.Kind() == reflect.Ptr {
					ptrToField.Set(reflect.ValueOf(fieldInstance))
				} else {
					ptrToField.Set(reflect.ValueOf(fieldInstance).Elem())
				}
			}
		}
	}
	finalType := instance.Type()
	if _, exists := c[finalType]; !exists {
		c[finalType] = make(map[string]binding)
	}
	c[finalType][regName] = binding{resolver: nil, instance: instance.Interface()}
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

	return nil, errors.ErrInvalidResolverSignature
}

func (c container) arguments(function any) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)

		if concrete, exist := c[abstraction][""]; exist {
			instance, _ := concrete.resolve(c)

			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, errors.ErrNoConcreteFound.WithMeta(abstraction.String())
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

// getTypeFromCache returns the reflect.Type for a type name, caching the result
func (c container) getTypeFromCache(typeName string, value any) reflect.Type {
	t, ok := reflectionCache.byTypeName[typeName]
	if ok {
		return t
	}
	t = reflect.TypeOf(value)
	reflectionCache.byTypeName[typeName] = t
	return t
}

// callLifecycleHook calls the given hook if implemented by the instance
func (c container) callLifecycleHook(instance any, hook string) {
	v := reflect.ValueOf(instance)
	method := v.MethodByName(hook)
	if method.IsValid() && method.Type().NumIn() == 0 {
		method.Call(nil)
	}
}

// RegisterInterface binds an interface type to a concrete implementation or factory.
func (c container) RegisterInterface(ifacePtr any, impl any, name string, singleton bool) error {
	ifaceType := c.getTypeFromCache("", ifacePtr)
	regName := name
	if regName == "" {
		regName = ifaceType.String()
	}
	if _, exist := c[ifaceType]; !exist {
		c[ifaceType] = make(map[string]binding)
	}
	var instance any
	var err error
	if singleton {
		// If impl is a function, invoke it to get the instance
		implType := reflect.TypeOf(impl)
		if implType.Kind() == reflect.Func {
			instance, err = c.invoke(impl)
			if err != nil {
				return err
			}
		} else {
			instance = impl
		}
		c[ifaceType][regName] = binding{resolver: impl, instance: instance}
	} else {
		c[ifaceType][regName] = binding{resolver: impl}
	}
	c.callLifecycleHook(instance, "OnRegister")
	return nil
}
