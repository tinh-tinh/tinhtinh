package validator

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"sync"
)

type Validator struct {
	cache sync.Map
}

func (v *Validator) Validate(obj any) error {
	val := reflect.ValueOf(obj)
	typ := val.Type()

	metaAny, ok := v.cache.Load(typ)
	if !ok {
		meta, err := v.compile(typ)
		if err != nil {
			return err
		}
		v.cache.Store(typ, meta)
		metaAny = meta
	}

	meta := metaAny.(*structMeta)

	return v.validateValue(val, meta)
}

func (v *Validator) validateNestedValue(fv reflect.Value, meta *structMeta) error {
	switch fv.Kind() {
	case reflect.Ptr:
		if fv.IsNil() {
			return nil
		}
		return v.validateValue(fv, meta)
	case reflect.Struct:
		return v.validateValue(fv, meta)

	case reflect.Slice, reflect.Array:
		for i := 0; i < fv.Len(); i++ {
			if err := v.validateNestedValue(fv.Index(i), meta); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Validator) validateValue(fv reflect.Value, meta *structMeta) error {
	if err := v.tryCustomScan(fv); err != nil {
		return err
	}

	if fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}

	var errMsg []string
	for _, f := range meta.fields {
		fieldVal := fv.FieldByIndex(f.index)
		fieldPath := f.name

		if isEmpty(fieldVal) && f.defaultValue.IsValid() && fieldVal.CanSet() {
			fieldVal.Set(f.defaultValue)
		}

		required := slices.IndexFunc(f.tags, func(v string) bool {
			return v == "required"
		})
		if required == -1 && isEmpty(fieldVal) {
			continue
		}

		for _, rule := range f.rules {
			if err := rule(fieldVal); err != nil {
				errMsg = append(errMsg, fieldPath+err.Error())
			}
		}

		if f.children != nil {
			if err := v.validateNestedValue(fieldVal, f.children); err != nil {
				errMsg = append(errMsg, err.Error())
			}
		}
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, "\n"))
	}
	return nil
}

func (v *Validator) tryCustomScan(fv reflect.Value) error {
	scannerType := reflect.TypeOf((*ScanFnc)(nil)).Elem()

	if fv.Type().Implements(scannerType) {
		res := fv.MethodByName("Scan").Call(nil)
		if !res[0].IsNil() {
			return res[0].Interface().(error)
		}
		return nil
	}

	if fv.CanAddr() {
		ptrFv := fv.Addr()
		if ptrFv.Type().Implements(scannerType) {
			res := ptrFv.MethodByName("Scan").Call(nil)
			if !res[0].IsNil() {
				return res[0].Interface().(error)
			}
			return nil
		}
	}
	return nil
}
