package valtruc

import "reflect"

func ptrValidatorWrapper(inner Validator, tag valTag) Validator {
	return func(ctx ValidationContext) (bool, error) {
		if ctx.FieldValue.Type().Kind() != reflect.Ptr {
			return inner(ctx)
		}

		z := reflect.Zero(ctx.FieldValue.Type())
		isZero := ctx.FieldValue.Interface() == z.Interface()
		if isZero {
			if tag.name == "required" {
				return false, NewValidationError(
					ctx,
					"the field pointer is nil and is required",
					RequiredIdentifier)
			} else {
				return true, nil
			}
		}

		ctx.FieldValue = ctx.FieldValue.Elem()
		return inner(ctx)
	}
}
