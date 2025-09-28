package valtruc

import (
	"fmt"
	"strconv"
)

const (
	MinInt64Identifier ValidatorIdentifier = "minInt64Identifier"
	MaxInt64Identifier ValidatorIdentifier = "maxInt64Identifier"
)

func minInt64(param string) Validator {
	minv, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min int64 %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Int()
		if value < minv {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("integer must be greater than %d", minv),
				MinInt64Identifier,
				param)
		}
		return true, nil
	}
}

func maxInt64(param string) Validator {
	maxv, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min int64 %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Int()
		if value > maxv {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("integer must be greater than %d", maxv),
				MaxInt64Identifier,
				param)
		}
		return true, nil
	}
}
