package validator

import (
	"errors"
	"reflect"
	"strconv"
)

type RuleFnc func(value reflect.Value) error

type RuleFactory func(param string) RuleFnc

var ruleRegister = map[string]RuleFactory{}

func RegisterRule(name string, factory RuleFactory) {
	ruleRegister[name] = factory
}

func init() {
	RegisterRule(tagRequired, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if isEmpty(value) {
				return errors.New(" is required")
			}
			return nil
		}
	})
	RegisterRule(tagIsAlpha, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isAlpha(value) {
				return errors.New(" is not a valid alpha")
			}
			return nil
		}
	})
	RegisterRule(tagIsAlphaNumeric, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isAlphanumeric(value) {
				return errors.New(" is not a valid alpha numeric")
			}
			return nil
		}
	})
	RegisterRule(tagIsEmail, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isEmail(value) {
				return errors.New(" is not a valid email")
			}
			return nil
		}
	})
	RegisterRule(tagIsUUID, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isUUID(value) {
				return errors.New(" is not a valid UUID")
			}
			return nil
		}
	})
	RegisterRule(tagIsObjectId, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isObjectId(value) {
				return errors.New(" is not a valid ObjectId")
			}
			return nil
		}
	})
	RegisterRule(tagIsStrongPassword, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isStrongPassword(value) {
				return errors.New(" is not a valid strong password")
			}
			return nil
		}
	})
	RegisterRule(tagIsInt, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isInt(value) {
				return errors.New(" is not a valid int")
			}
			return nil
		}
	})
	RegisterRule(tagIsFloat, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isFloat(value) {
				return errors.New(" is not a valid float")
			}
			return nil
		}
	})
	RegisterRule(tagIsNumber, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isInt(value) && !isFloat(value) {
				return errors.New(" is not a valid number")
			}
			return nil
		}
	})
	RegisterRule(tagIsDate, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isDate(value) {
				return errors.New(" is not a valid date")
			}
			return nil
		}
	})
	RegisterRule(tagIsDateString, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isDate(value) {
				return errors.New(" is not a valid date time")
			}
			return nil
		}
	})
	RegisterRule(tagIsBool, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			if !isBool(value) {
				return errors.New(" is not a valid bool")
			}
			return nil
		}
	})
	RegisterRule(tagNested, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			// no-op, handled separately
			return nil
		}
	})
	RegisterRule(tagMinLength, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			min, err := strconv.Atoi(param)
			if err != nil {
				return errors.New(" invalid minLength parameter")
			}
			if !minLength(value, min) {
				return errors.New(" minimum length is " + param)
			}
			return nil
		}
	})
	RegisterRule(tagMaxLength, func(param string) RuleFnc {
		return func(value reflect.Value) error {
			max, err := strconv.Atoi(param)
			if err != nil {
				return errors.New(" invalid maxLength parameter")
			}
			if !maxLength(value, max) {
				return errors.New(" maximum length is " + param)
			}
			return nil
		}
	})
}
