package valtruc

import (
	"fmt"
	"strconv"
)

const (
	MinSliceLengthIdentifier ValidatorIdentifier = "minSliceLengthIdentifier"
	MaxSliceLengthIdentifier ValidatorIdentifier = "maxSliceLengthIdentifier"
)

func minSliceLength(param string) Validator {
	minLen, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid max length string %s for slice", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		sliceLen := ctx.FieldValue.Len()
		if sliceLen < int(minLen) {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("the field required minimum length of %d", minLen),
				MinSliceLengthIdentifier,
				param)
		}
		return true, nil
	}
}

func maxSliceLength(param string) Validator {
	maxLen, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid max length string %s for slice", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		sliceLen := ctx.FieldValue.Len()
		if sliceLen > int(maxLen) {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("the field required minimum length of %d", maxLen),
				MaxSliceLengthIdentifier,
				param)
		}
		return true, nil
	}
}

func requiredSlice(_ string) Validator {
	return func(ctx ValidationContext) (bool, error) {
		isZero := ctx.FieldValue.IsNil()
		if isZero {
			return false, NewValidationError(
				ctx,
				"the field is required",
				RequiredIdentifier)
		}
		return true, nil
	}
}
