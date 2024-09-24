package utility

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const panicTextOfErrorRegistry = `
The error registry call sequence is incorrect.
It must start with AddErrorCode and end with NewError or WrapError
`

func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		errCodeMap:   make(map[int]*CustomError),
		errCodeSlice: make([]int, 0),
	}
}

type ErrorRegistry struct {
	errCodeMap   map[int]*CustomError // code:text
	errCodeSlice []int                // [code...]

	target *CustomError
}

func (r *ErrorRegistry) ShowErrors() {
	slices.SortFunc(r.errCodeSlice, func(a, b int) int {
		return a - b
	})

	texts := make([]string, 0, len(r.errCodeSlice)+2)

	texts = append(texts, "")
	for _, errCode := range r.errCodeSlice {
		text := fmt.Sprintf(" [ErrCode] %-8v %s", errCode, r.errCodeMap[errCode].Error())
		texts = append(texts, text)
	}
	texts = append(texts, "")

	fmt.Println(strings.Join(texts, "\n"))
}

// AddErrorCode ErrorCode is Must Field
func (r *ErrorRegistry) AddErrorCode(errCode int) *ErrorRegistry {
	if r.target != nil {
		panic(panicTextOfErrorRegistry)
	}

	_, exist := r.errCodeMap[errCode]
	if exist {
		panic("duplicated error code")
	}

	err := &CustomError{errCode: errCode}
	r.errCodeMap[errCode] = err
	r.errCodeSlice = append(r.errCodeSlice, errCode)

	r.target = err
	return r
}

// HttpStatus is Optional Field
func (r *ErrorRegistry) HttpStatus(httpStatus int) *ErrorRegistry {
	if r.target == nil {
		panic(panicTextOfErrorRegistry)
	}
	r.target.httpStatus = httpStatus
	return r
}

func (r *ErrorRegistry) NewError(description string) error {
	if r.target == nil {
		panic(panicTextOfErrorRegistry)
	}
	r.target.cause = errors.New(description)

	err := r.target
	r.target = nil
	return err
}

// WrapError
// 內部實作透過 fmt.Errorf 另外包裝 error,
// 模擬繼承的概念, 將 baseError 的 Optional Field 自動複製到新的 error
func (r *ErrorRegistry) WrapError(description string, baseError error) error {
	if r.target == nil {
		panic(panicTextOfErrorRegistry)
	}
	r.target.cause = fmt.Errorf("%v: %w", description, baseError)
	r.copyOptionalFieldsFrom(baseError)

	err := r.target
	r.target = nil
	return err
}

func (r *ErrorRegistry) copyOptionalFieldsFrom(baseError error) {
	if r.target.httpStatus == 0 {
		myErr, ok := UnwrapCustomError(baseError)
		if ok {
			r.target.httpStatus = myErr.httpStatus
		}
	}
}

//

func UnwrapCustomError(err error) (myErr *CustomError, ok bool) {
	if errors.As(err, &myErr) {
		return myErr, true
	}
	return nil, false
}

type CustomError struct {
	errCode    int
	httpStatus int
	cause      error
}

func (c *CustomError) Error() string {
	return c.cause.Error()
}

func (c *CustomError) ErrorCode() int {
	return c.errCode
}

func (c *CustomError) HttpStatus() int {
	return c.httpStatus
}

func (c *CustomError) Unwrap() error {
	return c.cause
}
