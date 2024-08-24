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
	min, err := strconv.ParseFloat(param, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min float64 %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Float()
		if value <= min {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("float must be greater than %f", min),
				MinFloat64Identifier,
				param)
		}
		return true, nil
	}
}

func maxFloat64(param string) Validator {
	max, err := strconv.ParseFloat(param, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min float64 %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Float()
		if value >= max {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("float must be greater than %f", max),
				MaxFloat64Identifier,
				param)
		}
		return true, nil
	}
}
