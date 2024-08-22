package valtruc_test

import (
	"strings"
	"testing"
	"valtruc"
)

func TestCore(t *testing.T) {
	t.Run("Should validate any struct without tags", func(t *testing.T) {
		type NoValtructValidations struct {
			Name string
		}

		ok, errs := valtruc.Validate(NoValtructValidations{})
		if !ok {
			t.Error("Validate should return that the struct is OK!")
		}
		if len(errs) > 0 {
			t.Error("Validate should return no errors")
		}
	})

	t.Run("Validate should only accept structs", func(t *testing.T) {
		defer func() {
			recover()
		}()
		valtruc.Validate(1)
		t.Error("Validate should panic when you pass something that is not a struct")
	})
}

func TestNestedStructs(t *testing.T) {
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
		ok, errs := valtruc.Validate(a{
			Name: "aab",
			Sub: b{
				Age: -1,
			},
		})
		if ok {
			t.Error("Validate should check substructs")
		}
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "integer must be greater than 0") {
			t.Error("The error returned should warn about the minimum value")
		}

		verr := errs[0].(valtruc.ValidationError)
		if verr.GetErrorCode() != valtruc.ErrCodeInt64Min {
			t.Error("The error returned should have ErrCodeInt64Min")
		}
		minValue, _ := verr.GetMetadataInt64("min")
		if minValue != 0 {
			t.Error("The required minimum value must be 0")
		}
	})

	t.Run("Should chceck substruct tags if substruct is not marked as required", func(t *testing.T) {
		ok, errs := valtruc.Validate(a{
			Name: "dd",
		})
		if ok {
			t.Error("Validate should check substructs")
		}
		if len(errs) == 0 {
			t.Error("Validate should return substruct errors")
		}
	})

	t.Run("Should complain that a requried substruct is missing", func(t *testing.T) {
		ok, errs := valtruc.Validate(c{
			Name: "i dont have a requried substruct",
		})
		if ok {
			t.Error("Validate should check substructs")
		}
		if len(errs) == 0 {
			t.Error("Validate should return substruct errors")
		}
	})
}

func TestRequired(t *testing.T) {
	t.Run("Empty struct with required fields should be not valid", func(t *testing.T) {
		type user struct {
			Name string `valtruc:"required"`
		}

		ok, errs := valtruc.Validate(user{})
		if ok {
			t.Error("Validate should check that name is not setted")
		}
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
	type user struct {
		Age int `valtruc:"min=18"`
	}

	t.Run("Struct with value inside limit should pass", func(t *testing.T) {
		ok, errs := valtruc.Validate(user{Age: 18})
		if ok {
			t.Error("Validate should check that int must have minimum")
		}
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "greater than 18") {
			t.Error("The error returned should warn about error")
		}
	})

	t.Run("Empty struct should not pass min", func(t *testing.T) {
		ok, errs := valtruc.Validate(user{})
		if ok {
			t.Error("Validate permit to use min and pass empty struct")
		}
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
	})
}

func TestStringValidators(t *testing.T) {
	type user struct {
		Name string `valtruc:"min=2"`
	}

	t.Run("Min string to reach minimum", func(t *testing.T) {
		ok, errs := valtruc.Validate(user{Name: "d"})
		if ok {
			t.Error("Validate should check that int must have minimum")
		}
		if len(errs) == 0 {
			t.Error("Validate should return at least one error")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
	})

	t.Run("Min string correct", func(t *testing.T) {
		ok, errs := valtruc.Validate(user{Name: "diego"})
		if !ok {
			t.Error("Validate should check that int is valid")
		}
		if len(errs) != 0 {
			t.Error("Validate should not return errors when its valid")
		}
	})

	type category struct {
		Name string `valtruc:"min=2, max=5"`
	}

	t.Run("String min should fail", func(t *testing.T) {
		ok, errs := valtruc.Validate(category{Name: "d"})
		if ok {
			t.Error("Validate should check that int must have minimum")
		}
		if len(errs) == 0 {
			t.Error("Validate should return at least one error (when checking minimum)")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
	})

	t.Run("String max should fail", func(t *testing.T) {
		ok, errs := valtruc.Validate(category{Name: "delgado"})
		if ok {
			t.Error("Validate should check that int must have maximum")
		}
		if len(errs) == 0 {
			t.Error("Validate should return at least one error (when checking maximum)")
		}
		if len(errs) > 1 {
			t.Error("Validate should return only one error")
		}
	})

	t.Run("String min max is correct", func(t *testing.T) {
		ok, errs := valtruc.Validate(category{Name: "diego"})
		if !ok {
			t.Error("Validate should check that int is valid")
		}
		if len(errs) != 0 {
			t.Error("Validate should not return errors when its valid")
		}
	})

	type tag struct {
		Name string `valtruc:"contains=pepo kawai, min=2"`
	}

	t.Run("String contains all ok", func(t *testing.T) {
		ok, errs := valtruc.Validate(tag{
			Name: "aqui viene el gran pepo kawai, una ranita muy simpatica",
		})
		if !ok {
			t.Error("Validate should check that string contains substring")
		}
		if len(errs) != 0 {
			t.Error("Validate should not return errors when its valid")
		}
	})

	t.Run("String contains failure", func(t *testing.T) {
		ok, errs := valtruc.Validate(tag{
			Name: "aqui viene el gran jamoncito, una ranita muy simpatica",
		})
		if ok {
			t.Error("Validate should check that string contains substring (pepo kawai)")
		}
		if len(errs) != 1 {
			t.Error("Validate should return only one error")
		}
		if !strings.Contains(errs[0].Error(), "must contain substring") {
			t.Error("The returned error should warn about the field must contain substring")
		}
	})
}
