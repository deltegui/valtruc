package valtruc

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidatorIdentifier string

type ValidationError struct {
	ctx        ValidationContext
	msg        string
	identifier ValidatorIdentifier
	param      string
}

func NewValidationError(ctx ValidationContext, msg string, identifier ValidatorIdentifier) ValidationError {
	return ValidationError{
		ctx:        ctx,
		msg:        msg,
		identifier: identifier,
		param:      "",
	}
}

func NewValidationErrorMeta(ctx ValidationContext, msg string, identifier ValidatorIdentifier, param string) ValidationError {
	return ValidationError{
		ctx:        ctx,
		msg:        msg,
		identifier: identifier,
		param:      param,
	}
}

func (verr ValidationError) GetStructName() string {
	return verr.ctx.StructType.Name()
}

func (verr ValidationError) GetFieldName() string {
	return verr.ctx.Field.Name
}

func (verr ValidationError) GetFieldTypeName() string {
	return verr.ctx.Field.Type.Name()
}

func (verr ValidationError) GetIdentifier() ValidatorIdentifier {
	return verr.identifier
}

func (verr ValidationError) GetFieldValue() string {
	switch verr.ctx.Field.Type.Kind() {
	case reflect.Int:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		return strconv.FormatInt(verr.ctx.FieldValue.Int(), 10)
	case reflect.Uint:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		return strconv.FormatUint(verr.ctx.FieldValue.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(verr.ctx.FieldValue.Bool())
	case reflect.String:
	default:
		return verr.ctx.FieldValue.String()
	}
	return verr.ctx.FieldValue.String()
}

func (verr ValidationError) GetParam() string {
	return verr.param
}

func (err ValidationError) Error() string {
	return fmt.Sprintf(
		"Validation error on struct '%s', field '%s' (%s) with value '%s': [%s] %s",
		err.GetStructName(),
		err.GetFieldName(),
		err.GetFieldTypeName(),
		err.GetFieldValue(),
		err.GetIdentifier(),
		err.msg)
}

func (err ValidationError) Format(str string) string {
	init := []rune(str)
	final := make([]rune, 0, len(str))

	for i := 0; i < len(init); i++ {
		c := init[i]
		if c != '$' {
			final = append(final, c)
			continue
		}

		remainingLen := len(init) - i
		if remainingLen >= 2 && str[i:i+3] == "${}" {
			value := err.GetParam()
			final = append(final, []rune(value)...)
			i += 2
			continue
		}

		final = append(final, '$')
	}

	return string(final)
}

type ValidationContext struct {
	StructType reflect.Type
	Field      reflect.StructField
	FieldIndex int
	FieldValue reflect.Value
}

type Validator func(ctx ValidationContext) (bool, error)
type ValidatorConstructor func(param string) Validator

type compiledValidation struct {
	validators []Validator
}

func (cValidation compiledValidation) validate(ctx ValidationContext) (bool, []error) {
	result := true
	errors := []error{}
	for _, validator := range cValidation.validators {
		ok, error := validator(ctx)
		if !ok {
			errors = append(errors, error)
		}
		result = result && ok
	}
	return result, errors
}

type Valtruc struct {
	compiled   map[reflect.Type]map[string]compiledValidation
	validators map[reflect.Kind]map[string]ValidatorConstructor
}

func New() Valtruc {
	return Valtruc{
		compiled:   map[reflect.Type]map[string]compiledValidation{},
		validators: builtInValidators,
	}
}

func (vt *Valtruc) AddValidator(forKind reflect.Kind, tagName string, constructor ValidatorConstructor) {
	vt.validators[forKind][tagName] = constructor
}

func (vt Valtruc) addCompilation(t reflect.Type, field string, value compiledValidation) {
	e, ok := vt.compiled[t]
	if ok {
		e[field] = value
		return
	}
	vt.compiled[t] = map[string]compiledValidation{
		field: value,
	}
}

func (vt Valtruc) Validate(target interface{}) []error {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	cc, ok := vt.compiled[t]
	if !ok {
		vt.compileStructValidation(t)
		cc = vt.compiled[t]
	}

	return vt.runValidations(t, v, cc)
}

func (vt Valtruc) runValidations(t reflect.Type, v reflect.Value, cc map[string]compiledValidation) []error {
	resultErrors := []error{}
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		ctx := ValidationContext{
			StructType: t,
			Field:      fieldType,
			FieldValue: fieldValue,
			FieldIndex: i,
		}

		validator := cc[fieldType.Name]
		validationResult, errors := validator.validate(ctx)
		if !validationResult {
			resultErrors = append(resultErrors, errors...)
		}

		if fieldType.Type.Kind() == reflect.Struct {
			errors := vt.runValidations(fieldType.Type, fieldValue, vt.compiled[fieldType.Type])
			resultErrors = append(resultErrors, errors...)
		}
	}
	return resultErrors
}

type valTag struct {
	structType reflect.Type
	field      reflect.StructField
	original   string
	name       string
	parameter  string
}

func (vt Valtruc) compileStructValidation(t reflect.Type) {
	if t.Kind() != reflect.Struct {
		panic("valtruc.Validate only accepts structs!")
	}
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		fieldType := t.Field(i)

		if fieldType.Type.Kind() == reflect.Struct {
			vt.compileStructValidation(fieldType.Type)
		}

		tag := fieldType.Tag
		val, ok := tag.Lookup("valtruc")
		if !ok {
			continue
		}

		tags := parseValtrucTag(val, fieldType, t)
		cc := vt.compile(tags)
		vt.addCompilation(t, fieldType.Name, cc)
	}
}

func parseValtrucTag(tag string, field reflect.StructField, structType reflect.Type) []valTag {
	tags := strings.Split(tag, ",")
	result := make([]valTag, len(tags))
	for i := 0; i < len(tags); i++ {
		t := strings.TrimSpace(tags[i])

		var name, param string
		startParamsIndex := strings.IndexRune(t, '=')
		if startParamsIndex != -1 {
			name = t[0:startParamsIndex]

			rest := t[startParamsIndex:]
			rest = strings.ReplaceAll(rest, "=", "")
			param = rest
		} else {
			name = t
		}

		result[i].structType = structType
		result[i].field = field
		result[i].original = t
		result[i].name = name
		result[i].parameter = param
	}

	return result
}

func (vt Valtruc) compile(tags []valTag) compiledValidation {
	result := compiledValidation{}

	for _, tag := range tags {
		validatorsForKind, ok := vt.validators[tag.field.Type.Kind()]
		if !ok {
			panic(fmt.Sprintf("valtruc: there is no validators for kind %s ", tag.field.Type.Kind()))
		}
		constructor, ok := validatorsForKind[tag.name]
		if !ok {
			panic(fmt.Sprintf("valtruc: validator with name %s not found for kind %s", tag.name, tag.field.Type.Kind()))
		}
		validator := constructor(tag.parameter)
		result.validators = append(result.validators, validator)
	}

	return result
}
