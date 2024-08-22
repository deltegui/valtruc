package valtruc

import (
	"fmt"
	"strconv"
)

func minInt64(param string) Validator {
	min, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min int64 %s", param))
	}
	meta := errMeta{
		"min": min,
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Int()
		if value <= min {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("integer must be greater than %d", min),
				ErrCodeInt64Min,
				meta)
		}
		return true, nil
	}
}

func maxInt64(param string) Validator {
	max, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min int64 %s", param))
	}
	meta := errMeta{
		"max": max,
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Int()
		if value >= max {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("integer must be greater than %d", max),
				ErrCodeInt64Max,
				meta)
		}
		return true, nil
	}
}
