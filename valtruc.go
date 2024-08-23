package valtruc

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type errMeta map[string]any

type ValidationError struct {
	ctx      ValidationContext
	msg      string
	code     ErrCode
	metadata errMeta
}

func NewValidationError(ctx ValidationContext, msg string, code ErrCode) ValidationError {
	return NewValidationErrorMeta(ctx, msg, code, errMeta{})
}
func NewValidationErrorMeta(ctx ValidationContext, msg string, code ErrCode, metadata errMeta) ValidationError {
	return ValidationError{
		ctx,
		msg,
		code,
		metadata,
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

func (verr ValidationError) GetErrorCode() ErrCode {
	return verr.code
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

func (verr ValidationError) GetMetadata(key string) (string, bool) {
	val, ok := verr.metadata[key]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

func (verr ValidationError) GetMetadataInt64(key string) (int64, bool) {
	val, ok := verr.metadata[key]
	if !ok {
		return 0, false
	}
	i, ok := val.(int64)
	return i, ok
}

func (err ValidationError) Error() string {
	return fmt.Sprintf(
		"Validation error on struct '%s', field '%s' (%s) with value '%s': [%d] %s",
		err.GetStructName(),
		err.GetFieldName(),
		err.GetFieldTypeName(),
		err.GetFieldValue(),
		err.code,
		err.msg)
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
	compiled map[reflect.Type]map[string]compiledValidation
}

func New() Valtruc {
	return Valtruc{map[reflect.Type]map[string]compiledValidation{}}
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

func (vt Valtruc) Validate(target interface{}) (bool, []error) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	cc, ok := vt.compiled[t]
	if !ok {
		vt.compileStructValidation(t)
		cc = vt.compiled[t]
	}

	return vt.runValidations(t, v, cc)
}

func (vt Valtruc) runValidations(t reflect.Type, v reflect.Value, cc map[string]compiledValidation) (bool, []error) {
	result := true
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
			result = false
		}

		if fieldType.Type.Kind() == reflect.Struct {
			ok, errors := vt.runValidations(fieldType.Type, fieldValue, vt.compiled[fieldType.Type])
			result = result && ok
			resultErrors = append(resultErrors, errors...)
		}
	}
	return result, resultErrors
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
		cc := compile(tags)
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

func compile(tags []valTag) compiledValidation {
	result := compiledValidation{}

	for _, tag := range tags {
		validatorsForKind, ok := validators[tag.field.Type.Kind()]
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
