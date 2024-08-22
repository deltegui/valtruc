package valtruc

import "reflect"

var intValidators map[string]ValidatorConstructor = map[string]ValidatorConstructor{
	"required": require,
	"min":      minInt64,
	"max":      maxInt64,
}

var stringValidators map[string]ValidatorConstructor = map[string]ValidatorConstructor{
	"required": require,
	"min":      minStringLength,
	"max":      maxStringLength,
	"contains": containsString,
}

var floatValidators map[string]ValidatorConstructor = map[string]ValidatorConstructor{
	"required": require,
}

var boolValidators map[string]ValidatorConstructor = map[string]ValidatorConstructor{
	"required": require,
}

var structValidators map[string]ValidatorConstructor = map[string]ValidatorConstructor{
	"required": require,
}

var validators map[reflect.Kind]map[string]ValidatorConstructor = map[reflect.Kind]map[string]ValidatorConstructor{
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
}

func AddValidator(forKind reflect.Kind, tagName string, constructor ValidatorConstructor) {
	// TODO IMPLEMENT
}
