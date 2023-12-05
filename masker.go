package x

import (
	"reflect"
)

type IMasker interface {
	Mask() string
}

type MaskReflect struct {
}

func (r MaskReflect) takeTyper(value reflect.Value) (ret reflect.Type) {
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return value.Type()
		}

		return value.Elem().Type()
	}

	valueType := value.Type()
	if valueType.Kind() == reflect.Ptr {
		return valueType.Elem()
	}

	return value.Type()
}

func (r MaskReflect) withMaskValue(valueField reflect.Value) (value reflect.Value) {
	if valueField.IsZero() {
		return valueField
	}

	var fieldType reflect.Type
	if valueField.Kind() == reflect.Ptr || valueField.Kind() == reflect.Interface {
		fieldType = valueField.Elem().Type()
	} else {
		fieldType = valueField.Type()
	}

	value = reflect.New(fieldType).Elem()
	if !fieldType.Implements(reflect.TypeOf((*IMasker)(nil)).Elem()) {
		return valueField
	}

	// 执行mask
	masker := valueField.Interface().(IMasker)
	maskString := masker.Mask()
	maskValue := reflect.ValueOf(maskString)

	if maskValue.CanConvert(fieldType) {
		value.Set(maskValue.Convert(fieldType))
		return value
	}

	value.Set(maskValue)
	return value
}

func (r MaskReflect) recursiveHandle(value reflect.Value) reflect.Value {
	if value.IsZero() {
		return value
	}

	typer := r.takeTyper(value)

	switch typer.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		return r.doMask(value.Interface())

	case reflect.String:
		return r.withMaskValue(value)
	}

	return value

}

func (r MaskReflect) cloneValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return reflect.New(value.Elem().Type())
	}

	return reflect.New(value.Type())
}

func (r MaskReflect) doStructMask(typer reflect.Type, val interface{}) reflect.Value {
	value := reflect.ValueOf(val)
	if value.IsZero() {
		return value
	}

	clone := r.cloneValue(value)
	if clone.Elem().Kind() == reflect.Struct && value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for i := 0; i < typer.NumField(); i++ {
		var (
			valueField reflect.Value
			cloneField reflect.Value
		)
		if value.Kind() == reflect.Ptr {
			valueField = value.Elem().Field(i)
		} else {
			valueField = value.Field(i)
		}

		cloneField = clone.Elem().Field(i)

		newValue := r.recursiveHandle(valueField)
		if newValue.CanConvert(cloneField.Type()) {
			cloneField.Set(newValue.Convert(cloneField.Type()))
			continue
		}

		valueKind, cloneKind := newValue.Kind(), cloneField.Kind()

		if valueKind == cloneKind {
			cloneField.Set(newValue)
			continue
		}

		if cloneField.Kind() == reflect.Struct && newValue.Kind() == reflect.Ptr {
			cloneField.Set(newValue.Elem())
		}

	}

	return clone
}

func (r MaskReflect) doSliceMask(typer reflect.Type, val interface{}) reflect.Value {
	value := reflect.ValueOf(val)
	clone := reflect.MakeSlice(typer, value.Len(), value.Len())

	for i := 0; i < value.Len(); i++ {
		cloneField := clone.Index(i)
		valuer := value.Index(i)

		cloneField.Set(
			r.recursiveHandle(valuer),
		)
	}

	return clone
}

func (r MaskReflect) doMappingMask(typer reflect.Type, val interface{}) reflect.Value {
	value := reflect.ValueOf(val)
	clone := reflect.MakeMap(typer)

	keys := value.MapKeys()
	for _, key := range keys {
		valuer := value.MapIndex(key)
		newValue := r.recursiveHandle(valuer)
		clone.SetMapIndex(key, newValue)
	}

	return clone
}

func (r MaskReflect) doMask(val interface{}) reflect.Value {
	var (
		typer reflect.Type
	)

	value := reflect.ValueOf(val)
	typer = r.takeTyper(value)

	switch typer.Kind() {
	case reflect.Struct:
		return r.doStructMask(typer, val)

	case reflect.Slice:
		return r.doSliceMask(typer, val)

	case reflect.Map:
		return r.doMappingMask(typer, val)
	case reflect.String:
		return r.withMaskValue(reflect.ValueOf(val))
	}

	return value
}

func (r MaskReflect) MakeMask(data interface{}) interface{} {
	return r.doMask(data).Interface()
}
