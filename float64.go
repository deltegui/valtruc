package valtruc

import (
	"fmt"
	"strconv"
)

const (
	MinFloat64Identifier ValidatorIdentifier = "minFloat64Identifier"
	MaxFloat64Identifier ValidatorIdentifier = "maxFloat64Identifier"
)

func minFloat64(param string) Validator {
	minv, err := strconv.ParseFloat(param, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min float64 %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Float()
		if value <= minv {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("float must be greater than %f", minv),
				MinFloat64Identifier,
				param)
		}
		return true, nil
	}
}

func maxFloat64(param string) Validator {
	maxv, err := strconv.ParseFloat(param, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min float64 %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Float()
		if value >= maxv {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("float must be greater than %f", maxv),
				MaxFloat64Identifier,
				param)
		}
		return true, nil
	}
}
