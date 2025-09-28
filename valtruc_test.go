package valtruc_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/deltegui/valtruc"
)

func TestCore(t *testing.T) {
	vt := valtruc.New()

	t.Run("Should validate any struct without tags", func(t *testing.T) {
		type NoValtructValidations struct {
			Name string
		}

		errs := vt.Validate(NoValtructValidations{})
		if errs != nil {
			t.Error("Validate should return no errors")
		}
	})

	t.Run("Validate should only accept structs", func(t *testing.T) {
		defer func() {
			_ = recover()
		}()
		vt.Validate(1)
		t.Error("Validate should panic when you pass something that is not a struct")
	})
}

func TestNestedStructsInsideSlices(t *testing.T) {
	vt := valtruc.New()

	type b struct {
		Age int32 `valtruc:"required, min=0, max=120"`
	}

	type a struct {
		Name string `valtruc:"min=2, max=255"`
		BS   []b
	}

	type c struct {
		Name string `valtruc:"required"`
		BS   [3]b
	}

	t.Run("Should check substructs inside slices", func(t *testing.T) {
		errs := vt.Validate(a{
			Name: "aab",
			BS: []b{
				{
					Age: -1,
				},
			},
		})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "integer must be greater than 0") {
			t.Error("The error returned should warn about the minimum value")
		}

		verr := valtruc.ValidationError{}
		ok := errors.As(errs[0], &verr)
		if !ok {
			t.Error("Expected errs[0] to be valtruc.ValidationError")
		}
		if verr.GetIdentifier() != valtruc.MinInt64Identifier {
			t.Error("The error returned should have MinInt64Identifier")
		}
		minValue := verr.GetParam()
		if minValue != "0" {
			t.Error("The required minimum value must be 0")
		}
	})

	t.Run("Should check substructs inside arrays", func(t *testing.T) {
		errs := vt.Validate(c{
			Name: "aab",
			BS: [3]b{
				{
					Age: -1,
				},
				{
					Age: 5,
				},
				{
					Age: -10,
				},
			},
		})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) != 2 {
			t.Error("Validate should return only two errors")
		}
		if !strings.Contains(errs[0].Error(), "integer must be greater than 0") {
			t.Error("The error returned should warn about the minimum value")
		}

		verr := valtruc.ValidationError{}
		ok := errors.As(errs[0], &verr)
		if !ok {
			t.Error("Expected errs[0] to be valtruc.ValidationError")
		}
		if verr.GetIdentifier() != valtruc.MinInt64Identifier {
			t.Error("The error returned should have MinInt64Identifier")
		}
		minValue := verr.GetParam()
		if minValue != "0" {
			t.Error("The required minimum value must be 0")
		}
	})
}

func TestNestedStructs(t *testing.T) {
	vt := valtruc.New()

	type b struct {
		Age int32 `valtruc:"required, min=0, max=120"`
	}

	type a struct {
		Name string `valtruc:"min=2, max=255"`
		Sub  b
	}

	type c struct {
		Name string `valtruc:"required"`
		Sub  b      `valtruc:"required"`
	}

	t.Run("Should check substructs", func(t *testing.T) {
		errs := vt.Validate(a{
			Name: "aab",
			Sub: b{
				Age: -1,
			},
		})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "integer must be greater than 0") {
			t.Error("The error returned should warn about the minimum value")
		}

		verr := valtruc.ValidationError{}
		ok := errors.As(errs[0], &verr)
		if !ok {
			t.Error("Expected errs[0] to be valtruc.ValidationError")
		}
		if verr.GetIdentifier() != valtruc.MinInt64Identifier {
			t.Error("The error returned should have MinInt64Identifier")
		}
		minValue := verr.GetParam()
		if minValue != "0" {
			t.Error("The required minimum value must be 0")
		}
	})

	t.Run("Should chceck substruct tags if substruct is not marked as required", func(t *testing.T) {
		errs := vt.Validate(a{
			Name: "dd",
		})
		if len(errs) == 0 {
			t.Error("Validate should return substruct errors")
		}
	})

	t.Run("Should complain that a requried substruct is missing", func(t *testing.T) {
		errs := vt.Validate(c{
			Name: "i dont have a requried substruct",
		})
		if len(errs) == 0 {
			t.Error("Validate should return substruct errors")
		}
	})
}

func TestRequired(t *testing.T) {
	vt := valtruc.New()

	t.Run("Empty struct with required fields should be not valid", func(t *testing.T) {
		type user struct {
			Name string `valtruc:"required"`
		}

		errs := vt.Validate(user{})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "field is required") {
			t.Error("The error returned should warn about the field is requried")
		}
	})
}

func TestMinInt(t *testing.T) {
	vt := valtruc.New()

	type user struct {
		Age int `valtruc:"min=18"`
	}

	t.Run("Struct with value inside limit should pass", func(t *testing.T) {
		errs := vt.Validate(user{Age: 18})
		if len(errs) > 0 {
			t.Error("Validate should return no errors")
		}
	})

	t.Run("Empty struct should not pass min", func(t *testing.T) {
		errs := vt.Validate(user{})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
	})
}

func TestStringValidators(t *testing.T) {
	vt := valtruc.New()

	type user struct {
		Name string `valtruc:"min=2"`
	}

	t.Run("Min string to reach minimum", func(t *testing.T) {
		errs := vt.Validate(user{Name: "d"})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
	})

	t.Run("Min string correct", func(t *testing.T) {
		errs := vt.Validate(user{Name: "diego"})
		if len(errs) != 0 {
			t.Error("Validate should not return errors when its valid")
		}
	})

	type category struct {
		Name string `valtruc:"min=2, max=5"`
	}

	t.Run("String min should fail", func(t *testing.T) {
		errs := vt.Validate(category{Name: "d"})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error (when checking minimum)")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
	})

	t.Run("String max should fail", func(t *testing.T) {
		errs := vt.Validate(category{Name: "delgado"})
		if len(errs) == 0 {
			t.Error("Validate should return at least one error (when checking maximum)")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
	})

	t.Run("String min max is correct", func(t *testing.T) {
		errs := vt.Validate(category{Name: "diego"})
		if len(errs) != 0 {
			t.Error("Validate should not return errors when its valid")
		}
	})

	type tag struct {
		Name string `valtruc:"contains=pepo kawai, min=2"`
	}

	t.Run("String contains all ok", func(t *testing.T) {
		errs := vt.Validate(tag{
			Name: "aqui viene el gran pepo kawai, una ranita muy simpatica",
		})
		if len(errs) != 0 {
			t.Error("Validate should not return errors when its valid")
		}
	})

	t.Run("String contains failure", func(t *testing.T) {
		errs := vt.Validate(tag{
			Name: "aqui viene el gran jamoncito, una ranita muy simpatica",
		})
		if len(errs) != 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "must contain substring") {
			t.Error("The returned error should warn about the field must contain substring")
		}
	})
}

func TestFormat(t *testing.T) {
	type tag struct {
		Name string `valtruc:"contains=pepo kawai, min=2"`
	}

	vt := valtruc.New()

	t.Run("Format should format user defined string", func(t *testing.T) {
		errs := vt.Validate(tag{Name: ""})
		if len(errs) == 1 {
			t.Error("Validate should return many error")
		}
		verr := valtruc.ValidationError{}
		ok := errors.As(errs[0], &verr)
		if !ok {
			t.Error("Expected errs[0] to be valtruc.ValidationError")
		}
		formatted := verr.Format("El nombre %s debe contener al menos la cadena: '${}'")
		if formatted != "El nombre %s debe contener al menos la cadena: 'pepo kawai'" {
			t.Error("The formatted error should warn about the field must contain substring pepo kawai")
		}
	})
}

func TestCanAddCustomValidators(t *testing.T) {
	type tag struct {
		Name string `valtruc:"reverse=iawak, min=2"`
	}

	vt := valtruc.New()

	const identifier valtruc.ValidatorIdentifier = "reverseStringIdentifier"

	reverse := func(param string) valtruc.Validator {
		return func(ctx valtruc.ValidationContext) (bool, error) {
			str := ctx.FieldValue.String()
			i := 0              // str
			j := len(param) - 1 // param
			for i < len(str) && j >= 0 {
				if param[j] != str[i] {
					return false, valtruc.NewValidationErrorMeta(
						ctx,
						"item is not reverse",
						identifier,
						param)
				}
				j--
				i++
			}
			return true, nil
		}
	}

	vt.AddValidator(reflect.String, "reverse", reverse)

	t.Run("Should check custom validators", func(t *testing.T) {
		errs := vt.Validate(tag{Name: "kawasaki"})
		if errs == nil {
			t.Error("Validate should return one error")
		}
		if len(errs) != 1 {
			t.Error("Validate should return one error")
		}
		verr := valtruc.ValidationError{}
		ok := errors.As(errs[0], &verr)
		if !ok {
			t.Error("Expected errs[0] to be valtruc.ValidationError")
		}
		if verr.GetIdentifier() != identifier {
			t.Error("The error code must be correct")
		}
		formatted := verr.Format("El elemento debe ser la cadena revertida de: '${}'")
		if formatted != "El elemento debe ser la cadena revertida de: 'iawak'" {
			t.Error("The formatted error should tell about the reversed string")
		}
	})
}

func TestCompleteValidation(t *testing.T) {
	type User struct {
		Name             string `valtruc:"min=3, max=10, required"`
		Password         string `valtruc:"min=3, max=255, required"`
		AcceptConditions bool   `valtruc:"mustBeTrue"`
		Email            string `valtruc:"min=3, max=255, required"`
	}

	vt := valtruc.New()

	user := User{
		Name:             "abcdefghijklmnopqrst",
		Password:         "b",
		Email:            "c",
		AcceptConditions: false,
	}

	errs := vt.Validate(user)

	if len(errs) == 0 {
		t.Error("Must be at least one error")
	}
}

func TestShouldHandlePointers(t *testing.T) {
	type UserFilter struct {
		Name             *string `valtruc:"min=3, max=10, required"`
		Password         *string `valtruc:"min=3, max=255, required"`
		AcceptConditions *bool   `valtruc:"mustBeTrue"`
		Email            *string `valtruc:"min=3, max=255, required"`
	}

	t.Run("Struct with pointers all ok", func(t *testing.T) {
		vt := valtruc.New()

		var (
			name             string = "diego"
			password                = "ddb"
			email                   = "diego@deltegui.com"
			acceptConditions bool   = true
		)
		user := UserFilter{
			Name:             &name,
			Password:         &password,
			Email:            &email,
			AcceptConditions: &acceptConditions,
		}

		errs := vt.Validate(user)

		if len(errs) != 0 {
			t.Error("Should not be errors")
		}
	})

	t.Run("Struct with nil pointers should be ok when is not required", func(t *testing.T) {
		type UserFilter struct {
			Name             *string `valtruc:"min=3, max=10"`
			Password         *string `valtruc:"min=3, max=255"`
			AcceptConditions *bool   `valtruc:"mustBeTrue"`
			Email            *string `valtruc:"min=3, max=255"`
		}

		vt := valtruc.New()

		errs := vt.Validate(UserFilter{})

		if len(errs) != 0 {
			t.Error("Should not be errors")
		}
	})

	t.Run("Struct with nil pointers should fail when is required", func(t *testing.T) {
		type UserFilter struct {
			Name             *string `valtruc:"min=3, max=10, required"`
			Password         *string `valtruc:"min=3, max=255"`
			AcceptConditions *bool   `valtruc:"mustBeTrue"`
			Email            *string `valtruc:"min=3, max=255"`
		}

		vt := valtruc.New()

		errs := vt.Validate(UserFilter{})

		if len(errs) != 1 {
			t.Error("Should be one error")
		}
	})

	t.Run("Old validators should still run with pointers", func(t *testing.T) {
		type UserFilter struct {
			Name             *string `valtruc:"min=3, max=10"`
			Password         *string `valtruc:"min=3, max=255"`
			AcceptConditions *bool   `valtruc:"mustBeTrue"`
			Email            *string `valtruc:"min=3, max=255"`
		}

		vt := valtruc.New()

		var (
			name string = "d"
		)
		errs := vt.Validate(UserFilter{Name: &name})

		if len(errs) != 1 {
			t.Error("Should be one error from string")
		}
	})
}

func TestSliceValidators(t *testing.T) {
	type UserFilter struct {
		Roles []int `valtruc:"min=3, max=10, required"`
	}

	t.Run("Struct with slice ok", func(t *testing.T) {
		vt := valtruc.New()

		roles := []int{1, 2, 3}
		user := UserFilter{
			Roles: roles,
		}

		errs := vt.Validate(user)

		if len(errs) != 0 {
			t.Error("Should not be errors")
		}
	})
}
