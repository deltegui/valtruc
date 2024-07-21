package valtruc

import (
	"fmt"
	"strconv"
	"strings"
)

func minStringLength(param string) Validator {
	min, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min length string %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.String()
		if len(value) < int(min) {
			return false, NewValidationError(ctx, fmt.Sprintf("the field required minimum length of %d", min))
		}
		return true, nil
	}
}

func maxStringLength(param string) Validator {
	max, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid max length string %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.String()
		if len(value) > int(max) {
			return false, NewValidationError(ctx, fmt.Sprintf("the field required maximum length of %d", max))
		}
		return true, nil
	}
}

func containsString(param string) Validator {
	if len(param) == 0 {
		panic("string contains must have a parameter telling what contains")
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.String()
		if !strings.Contains(value, param) {
			return false, NewValidationError(ctx, fmt.Sprintf("the field must contain substring %s", param))
		}
		return true, nil
	}
}
