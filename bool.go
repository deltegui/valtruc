package valtruc

const (
	MustBeTrueBoolIdentifier  ValidatorIdentifier = "mustBeTrueBoolIdentifier"
	MustBeFalseBoolIdentifier ValidatorIdentifier = "mustBeFalseBoolIdentifier"
)

func mustBeTrue(_ string) Validator {
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Bool()
		if !value {
			return false, NewValidationError(
				ctx,
				"bool must be true",
				MustBeTrueBoolIdentifier)
		}
		return true, nil
	}
}

func mustBeFalse(_ string) Validator {
	return func(ctx ValidationContext) (bool, error) {
		value := ctx.FieldValue.Bool()
		if value {
			return false, NewValidationError(
				ctx,
				"bool must be false",
				MustBeFalseBoolIdentifier)
		}
		return true, nil
	}
}
