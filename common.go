package valtruc

import (
	"reflect"
)

func require(param string) Validator {
	return func(ctx ValidationContext) (bool, error) {
		z := reflect.Zero(ctx.FieldValue.Type())
		isZero := ctx.FieldValue.Interface() == z.Interface()
		if isZero {
			return false, NewValidationError(ctx, "the field is required")
		}
		return true, nil
	}
}
