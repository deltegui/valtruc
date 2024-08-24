# Valtruc: structure validator for Go

Valtruc is a simple, zero-dependency structure validator library for Go

## Installation
```
go install github.com/deltegui/valtruc@latest
```

## Basic usage
First add valtruc tags to your structs

```
type User struct {
    Name string `valtruc:"min=3, max=255, required"`
    Password string `valtruc:"min=3, max=255, required"`
    AcceptConditions bool `valtruc:"mustBeTrue"`
    Email string `valtruc:"min=3, max=255, required"`
}
```

then create and validate your struct

```
user := User{
    Name:             "a",
    Password:         "b",
    Email:            "c",
    AcceptConditions: false,
}

errs := vt.Validate(user)
```

The returned `errs` is a `error` array. You can iterate over it and print the error:

```
for _, err := range errs {
    fmt.Println(err)
}
```

This code will output:

```
Validation error on struct 'User', field 'Name' (string) with value 'abcdefghijklmnopqrst': [maxStringLengthIdentifier] the field required maximum length of 10
Validation error on struct 'User', field 'Password' (string) with value 'b': [minStringLengthIdentifier] the field required minimum length of 3
Validation error on struct 'User', field 'AcceptConditions' (bool) with value 'false': [mustBeTrueBoolIdentifier] bool must be true
Validation error on struct 'User', field 'Email' (string) with value 'c': [minStringLengthIdentifier] the field required minimum length of 3
```

## Error API
You can transform the returned `error` to `valtruc.ValidationError` type to access all validation error information. The available methods in `ValidationError` are:

* `GetStructName() string`: Gets the struct name (eg. `User`)
* `GetFieldName() string`: Get the validated field name (eg. `Email`)
* `GetFieldTypeName() string`: Get the validated field name (eg. `string`)
* `GetIdentifier() valtruc.ValidatorIdentifier`: Gets an validator identifier (type alias for string). You can use this to programmatically check the error type:
```
if (err.GetIdentifier() == valtruc.MinFloat64Identifier) {
    // handle error knowing is MinFloat64
}
```
* `GetFieldValue() string`: Get field value as string (eg. `10`)
* `GetParam() string`: Get validator param (if you have used min validator `min=2` the returned string is `2`)
* `Error() string`
* `Format(str string) string`: Formats the error. You can use `${}` placeholder to show the param value. (eg. `Format("Must be minimum of ${}")` will output `Must be minimum of 2`)

## Create your own validators
```
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
```

usage:

```
type tag struct {
    Name string `valtruc:"reverse=iawak, min=2"`
}
```
