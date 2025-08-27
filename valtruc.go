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

func (err ValidationError) Path() []string {
	return err.ctx.Path
}

func NewValidationError(
	ctx ValidationContext,
	msg string,
	identifier ValidatorIdentifier,
) ValidationError {
	return ValidationError{
		ctx:        ctx,
		msg:        msg,
		identifier: identifier,
		param:      "",
	}
}

func NewValidationErrorMeta(
	ctx ValidationContext,
	msg string,
	identifier ValidatorIdentifier,
	param string,
) ValidationError {
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

func (verr ValidationError) Error() string {
	return fmt.Sprintf(
		"Validation error on struct '%s', field '%s' (%s) with value '%s': [%s] %s",
		verr.GetStructName(),
		verr.GetFieldName(),
		verr.GetFieldTypeName(),
		verr.GetFieldValue(),
		verr.GetIdentifier(),
		verr.msg)
}

func FormatWithParam(str, param string) string {
	init := []rune(str)
	final := make([]rune, 0, len(str))

	for i := 0; i < len(init); i++ {
		c := init[i]
		if c != '$' {
			final = append(final, c)
			continue
		}

		remainingLen := len(init) - i
		if remainingLen >= 2 && string(init[i:i+3]) == "${}" {
			final = append(final, []rune(param)...)
			i += 2
			continue
		}

		final = append(final, '$')
	}

	return string(final)
}

func (verr ValidationError) Format(str string) string {
	return FormatWithParam(str, verr.GetParam())
}

type ValidationContext struct {
	StructType reflect.Type
	Field      reflect.StructField
	FieldIndex int
	FieldValue reflect.Value
	Path       []string
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
		ok, err := validator(ctx)
		if !ok {
			errors = append(errors, err)
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
		validators: createValidators(),
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

	errs := vt.runValidations(t, v, cc, []string{})
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (vt Valtruc) runValidations(t reflect.Type, v reflect.Value, cc map[string]compiledValidation, path []string) []error {
	resultErrors := []error{}
	numFields := t.NumField()
	for i := range numFields {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		ctx := ValidationContext{
			StructType: t,
			Field:      fieldType,
			FieldValue: fieldValue,
			FieldIndex: i,
			Path:       path,
		}

		validator := cc[fieldType.Name]
		validationResult, errors := validator.validate(ctx)
		if !validationResult {
			resultErrors = append(resultErrors, errors...)
		}

		if fieldType.Type.Kind() == reflect.Struct {
			subpath := append(path, fieldType.Name)
			validationErrors := vt.runValidations(fieldType.Type, fieldValue, vt.compiled[fieldType.Type], subpath)
			resultErrors = append(resultErrors, validationErrors...)
		}
		if fieldType.Type.Kind() == reflect.Array || fieldType.Type.Kind() == reflect.Slice {
			v := fieldValue
			for j := 0; j < v.Len(); j++ {
				indexed := v.Index(j)
				if indexed.Type().Kind() == reflect.Struct {
					subpath := append(path, fmt.Sprintf("%s[%d]", fieldType.Name, j))
					validationErrors := vt.runValidations(indexed.Type(), indexed, vt.compiled[indexed.Type()], subpath)
					resultErrors = append(resultErrors, validationErrors...)
				}
			}
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
	for i := range numFields {
		fieldType := t.Field(i)

		if fieldType.Type.Kind() == reflect.Struct {
			vt.compileStructValidation(fieldType.Type)
		}
		if fieldType.Type.Kind() == reflect.Array || fieldType.Type.Kind() == reflect.Slice {
			underlyingType := fieldType.Type.Elem()
			if underlyingType.Kind() == reflect.Struct {
				vt.compileStructValidation(underlyingType)
			}
		}

		tag := fieldType.Tag
		val, ok := tag.Lookup("valtruc")
		if !ok {
			continue
		}

		tags := parseValtrucTag(val, fieldType, t)
		cc := vt.compile(tags, fieldType)
		vt.addCompilation(t, fieldType.Name, cc)
	}
}

func parseValtrucTag(tag string, field reflect.StructField, structType reflect.Type) []valTag {
	tags := strings.Split(tag, ",")
	result := make([]valTag, len(tags))
	for i := range len(tags) {
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

func (vt Valtruc) compile(tags []valTag, field reflect.StructField) compiledValidation {
	result := compiledValidation{}

	kind := field.Type.Kind()
	isPtr := false

	if kind == reflect.Ptr {
		kind = field.Type.Elem().Kind()
		isPtr = true
	}

	for _, tag := range tags {
		validatorsForKind, ok := vt.validators[kind]
		if !ok {
			panic(fmt.Sprintf("valtruc: there is no validators for kind %s ", kind))
		}
		constructor, ok := validatorsForKind[tag.name]
		if !ok {
			panic(fmt.Sprintf("valtruc: validator with name %s not found for kind %s", tag.name, kind))
		}
		validator := constructor(tag.parameter)
		if isPtr {
			validator = ptrValidatorWrapper(validator, tag)
		}
		result.validators = append(result.validators, validator)
	}

	return result
}
