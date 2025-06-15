package valtruc

import "reflect"

func createValidators() map[reflect.Kind]map[string]ValidatorConstructor {

	var intValidators = map[string]ValidatorConstructor{
		"required": require,
		"min":      minInt64,
		"max":      maxInt64,
	}

	var stringValidators = map[string]ValidatorConstructor{
		"required": require,
		"min":      minStringLength,
		"max":      maxStringLength,
		"contains": containsString,
	}

	var floatValidators = map[string]ValidatorConstructor{
		"required": require,
		"min":      minFloat64,
		"max":      maxFloat64,
	}

	var boolValidators = map[string]ValidatorConstructor{
		"required":    require,
		"mustBeTrue":  mustBeTrue,
		"mustBeFalse": mustBeFalse,
	}

	var structValidators = map[string]ValidatorConstructor{
		"required": require,
	}

	var sliceValidators = map[string]ValidatorConstructor{
		"required": requiredSlice,
		"max":      maxSliceLength,
		"min":      minSliceLength,
	}

	return map[reflect.Kind]map[string]ValidatorConstructor{
		reflect.String:  stringValidators,
		reflect.Int:     intValidators,
		reflect.Int16:   intValidators,
		reflect.Int32:   intValidators,
		reflect.Int64:   intValidators,
		reflect.Uint:    intValidators,
		reflect.Uint16:  intValidators,
		reflect.Uint32:  intValidators,
		reflect.Uint64:  intValidators,
		reflect.Float32: floatValidators,
		reflect.Float64: floatValidators,
		reflect.Bool:    boolValidators,
		reflect.Struct:  structValidators,
		reflect.Slice:   sliceValidators,
	}
}
