package valtruc

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MinStringLengthIdentifier ValidatorIdentifier = "minStringLengthIdentifier"
	MaxStringLengthIdentifier ValidatorIdentifier = "maxStringLengthIdentifier"
	ContainsStringIdentifier  ValidatorIdentifier = "containsStringIdentifier"
)

func minStringLength(param string) Validator {
	minv, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid min length string %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.String()
		if len(value) < int(minv) {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("the field required minimum length of %d", minv),
				MinStringLengthIdentifier,
				param)
		}
		return true, nil
	}
}

func maxStringLength(param string) Validator {
	maxv, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid max length string %s", param))
	}
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.String()
		if len(value) > int(maxv) {
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("the field required maximum length of %d", maxv),
				MaxStringLengthIdentifier,
				param)
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
			return false, NewValidationErrorMeta(
				ctx,
				fmt.Sprintf("the field must contain substring %s", param),
				ContainsStringIdentifier,
				param)
		}
		return true, nil
	}
}
