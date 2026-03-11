package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type fieldMeta struct {
	index        []int
	name         string
	rules        []RuleFnc
	tags         []string
	children     *structMeta
	defaultValue reflect.Value
}

type structMeta struct {
	fields []fieldMeta
}

func (v *Validator) compile(t reflect.Type) (*structMeta, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("only struct supported")
	}

	meta := &structMeta{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagVal := field.Tag.Get("validate")
		if tagVal == "" {
			continue
		}

		rules, err := parseTag(tagVal)
		if err != nil {
			return nil, err
		}

		fieldType := field.Type
		children, err := v.compileNested(fieldType)
		if err != nil {
			return nil, err
		}

		if rules == nil && children == nil {
			continue
		}

		fMeta := fieldMeta{
			index:    field.Index,
			name:     field.Name,
			rules:    rules,
			tags:     strings.Split(tagVal, ","),
			children: children,
		}

		defaultTag := field.Tag.Get("default")
		if defaultTag != "" {
			defVal, err := parseDefaultValue(field.Type, defaultTag)
			if err != nil {
				return nil, err
			}
			fMeta.defaultValue = defVal
		}

		meta.fields = append(meta.fields, fMeta)
	}

	return meta, nil
}

func parseTag(tag string) ([]RuleFnc, error) {
	parts := strings.Split(tag, ",")
	rules := make([]RuleFnc, 0, len(parts))
	for _, part := range parts {
		name, param, _ := strings.Cut(part, "=")
		if factory, ok := ruleRegister[name]; ok {
			rule := factory(param)
			rules = append(rules, rule)
		} else {
			return nil, errors.New("unknown validator: " + name)
		}
	}

	return rules, nil
}

func (v *Validator) compileNested(t reflect.Type) (*structMeta, error) {
	switch t.Kind() {
	case reflect.Ptr:
		return v.compileNested(t.Elem())

	case reflect.Struct:
		if isPrimitiveStruct(t) {
			return nil, nil
		}
		return v.compile(t)

	case reflect.Slice, reflect.Array:
		elem := t.Elem()
		return v.compileNested(elem)
	}
	return nil, nil
}

func isPrimitiveStruct(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

func parseDefaultValue(t reflect.Type, val string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(val), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid int value '%s': %w", val, err)
		}
		out := reflect.New(t).Elem()
		out.SetInt(i)
		return out, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid uint value '%s': %w", val, err)
		}
		out := reflect.New(t).Elem()
		out.SetUint(ui)
		return out, nil

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid float value '%s': %w", val, err)
		}
		out := reflect.New(t).Elem()
		out.SetFloat(f)
		return out, nil

	case reflect.Bool:
		b, err := strconv.ParseBool(strings.ToLower(val))
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid bool value '%s': %w", val, err)
		}
		return reflect.ValueOf(b), nil
	}

	return reflect.Value{}, fmt.Errorf("unsupported default value type: %s", t.Kind())
}
