package valtruc

import "reflect"

const (
	RequiredIdentifier ValidatorIdentifier = "requiredIdentifier"
)

func require(_ string) Validator {
	return func(ctx ValidationContext) (bool, error) {
		z := reflect.Zero(ctx.FieldValue.Type())
		isZero := ctx.FieldValue.Interface() == z.Interface()
		if isZero {
			return false, NewValidationError(
				ctx,
				"the field is required",
				RequiredIdentifier)
		}
		return true, nil
	}
}
