package valtruc

type ErrCode int

const (
	// Common
	ErrCodeRequired ErrCode = 0

	// String
	ErrCodeStringMinLength ErrCode = 100
	ErrCodeStringMaxLength ErrCode = 101
	ErrCodeStringContains  ErrCode = 102

	// Int64
	ErrCodeInt64Min ErrCode = 200
	ErrCodeInt64Max ErrCode = 201
)
