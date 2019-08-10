package convertstruct

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Converter interface {
		Convert(value interface{}) error
	}

	converterImpl struct {
		original reflect.Value
	}
)

func NewConverter(original interface{}) Converter {
	return &converterImpl{
		original: reflect.ValueOf(original),
	}
}

func (c *converterImpl) Convert(destination interface{}) error {
	value := reflect.ValueOf(destination)

	if value.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, expected %#v", destination)
	}

	return c.convert([]string{}, c.original, value)
}

func (c *converterImpl) convert(paths []string, source reflect.Value, dest reflect.Value) error {
	if !dest.IsValid() {
		return nil
	}
	
	if c.isNilValue(source) {
		return nil
	}

	rootSource := c.getRootValue(source)
	rootDest := c.getRootValue(dest)

	if rootSource.Type().Kind() == rootDest.Type().Kind() {
		return c.convertSameKind(paths, rootSource, rootDest)
	}

	return nil
}

func (c *converterImpl) isNilValue(value reflect.Value) bool {
	switch value.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func (c *converterImpl) getRootValue(value reflect.Value) reflect.Value {
	switch value.Type().Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
			value.Elem().Set(reflect.Zero(value.Type().Elem()))
		}
		return c.getRootValue(value.Elem())
	default:
		return value
	}
}

func (c *converterImpl) convertSameKind(paths []string, source reflect.Value, dest reflect.Value) error {
	switch source.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return c.convertSlices(paths, source, dest)
	case reflect.Struct:
		return c.convertStructs(paths, source, dest)
	case reflect.Interface:
		return c.convertInterfaces(paths, source, dest)
	default:
		dest.Set(source)
		return nil
	}
}

func (c *converterImpl) convertSlices(paths []string, source reflect.Value, dest reflect.Value) error {
	sourceType := source.Type().Elem()
	destType := dest.Type().Elem()

	if !sourceType.AssignableTo(destType) && !sourceType.ConvertibleTo(destType) {
		return fmt.Errorf(`field "%s" as slice with type %s is not convertible to slice with type %s`, strings.Join(paths, "."), sourceType.String(), destType.String())
	}

	dest.Set(reflect.Zero(dest.Type()))
	for i := 0; i < source.Len(); i++ {
		if sourceType.Kind() != destType.Kind() {
			dest.Set(reflect.Append(dest, c.convertToValue(source.Index(i), destType)))
		} else {
			dest.Set(reflect.Append(dest, source.Index(i)))
		}
	}

	return nil
}

func (c *converterImpl) convertStructs(paths []string, source reflect.Value, dest reflect.Value) error {
	sourceType := source.Type()
	destType := dest.Type()

	if sourceType.AssignableTo(destType) {
		dest.Set(source)
		return nil
	}

	for i := 0; i< source.NumField(); i++ {
		fieldName := sourceType.Field(i).Name

		if _, ok := destType.FieldByName(fieldName); !ok {
			continue
		}

		err := c.convert(append(paths, fieldName), source.Field(i), dest.FieldByName(fieldName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *converterImpl) convertInterfaces(paths []string, source reflect.Value, dest reflect.Value) error {
	sourceType := source.Type()
	destType := dest.Type()

	if c.isCoreInterface(dest) || sourceType.AssignableTo(destType) {
		dest.Set(source)
		return nil
	}

	return nil
}

func (c *converterImpl) isCoreInterface(value reflect.Value) bool {
	return value.Type().String() == "interface {}"
}

func (c *converterImpl) convertToValue(source reflect.Value, destType reflect.Type) reflect.Value {
	return source.Convert(destType)
}