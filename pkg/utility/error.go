package utility

import (
	"errors"
	"fmt"
	"net/http"
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
	slices.SortFunc(r.errCodeSlice, func(a, b int) int { return a - b })

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
// The correct call sequence starts with AddErrorCode and ends with NewError or WrapError.
//
// Example:
//
//	ErrInvalidParam = ErrorRegistry.
//		AddErrorCode(4000).
//		AddHttpStatus(http.StatusBadRequest).
//		NewError("invalid parameter")
//
//	ErrInvalidUsername = ErrorRegistry.
//		AddErrorCode(6000).
//		WrapError("username must be having a upper letter", ErrInvalidParam)
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

func (r *ErrorRegistry) NewError(description string) error {
	if r.target == nil {
		panic(panicTextOfErrorRegistry)
	}
	r.target.cause = errors.New(description)

	err := r.target
	r.target = nil
	return err
}

// WrapError 透過 fmt.Errorf 包裝 error
func (r *ErrorRegistry) WrapError(description string, baseError error) error {
	if r.target == nil {
		panic(panicTextOfErrorRegistry)
	}
	r.target.cause = fmt.Errorf("%v: %w", description, baseError)
	r.target.copyFrom(baseError)

	err := r.target
	r.target = nil
	return err
}

func (r *ErrorRegistry) AddHttpStatus(httpStatus int) *ErrorRegistry {
	if r.target == nil {
		panic(panicTextOfErrorRegistry)
	}
	r.target.httpStatus = httpStatus
	return r
}

//

var ErrUnknown = &CustomError{
	cause:      errors.New("unknown error"),
	errCode:    -1,
	httpStatus: http.StatusInternalServerError,
}

func UnwrapCustomError(err error) (myErr *CustomError, ok bool) {
	if errors.As(err, &myErr) {
		return myErr, true
	}
	return ErrUnknown, false
}

type CustomError struct {
	cause   error
	errCode int

	// Optional Field
	httpStatus int
}

func (c *CustomError) Error() string {
	return c.cause.Error()
}

// ErrorCode 必須有明確定義，才可視為系統規範內的錯誤, 否則視為 ErrUnknown
func (c *CustomError) ErrorCode() int {
	return c.errCode
}

func (c *CustomError) Unwrap() error {
	return c.cause
}

// 將 baseErr 的 Optional Field 進行複製, 想模擬繼承的概念
func (c *CustomError) copyFrom(baseErr error) {
	Err, _ := UnwrapCustomError(baseErr)
	if c.httpStatus == 0 {
		c.httpStatus = Err.httpStatus
	}
}

func (c *CustomError) HttpStatus() int {
	return c.httpStatus
}
