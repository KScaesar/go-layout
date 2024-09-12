package utility

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		errCodeMap:   make(map[string]error),
		errCodeSlice: make([]string, 0),
	}
}

type ErrorRegistry struct {
	errCodeMap   map[string]error
	errCodeSlice []string

	target *CustomError
}

func (r *ErrorRegistry) ShowErrors() {
	slices.SortFunc(r.errCodeSlice, strings.Compare)

	texts := make([]string, 0, len(r.errCodeSlice)+2)

	texts = append(texts, "")
	for _, errCode := range r.errCodeSlice {
		text := fmt.Sprintf(" [ErrCode] %-8s: %s", errCode, r.errCodeMap[errCode].Error())
		texts = append(texts, text)
	}
	texts = append(texts, "")

	fmt.Println(strings.Join(texts, "\n"))
}

func (r *ErrorRegistry) AddError(errCode int, err error) error {
	return r.AddErrorCode(errCode).
		Description(err.Error()).
		NewError()
}

func (r *ErrorRegistry) AddErrorCode(errCode int) *ErrorRegistry {
	if r.target != nil {
		panic("The error registry call sequence is incorrect. It must start with AddErrorCode and end with NewError.")
	}

	myCode := strconv.Itoa(errCode)
	_, exist := r.errCodeMap[myCode]
	if exist {
		panic("duplicated custom code")
	}

	err := &CustomError{errCode: errCode}
	r.errCodeMap[myCode] = err
	r.errCodeSlice = append(r.errCodeSlice, myCode)

	r.target = err
	return r
}

func (r *ErrorRegistry) HttpStatus(httpStatus int) *ErrorRegistry {
	if r.target == nil {
		panic("The error registry call sequence is incorrect. It must start with AddErrorCode and end with NewError.")
	}
	r.target.httpStatus = httpStatus
	return r
}

func (r *ErrorRegistry) Description(description string) *ErrorRegistry {
	if r.target == nil {
		panic("The error registry call sequence is incorrect. It must start with AddErrorCode and end with NewError.")
	}
	r.target.description = description
	return r
}

func (r *ErrorRegistry) NewError() error {
	if r.target == nil {
		panic("The error registry call sequence is incorrect. It must start with AddErrorCode and end with NewError.")
	}
	err := r.target
	r.target = nil
	return err
}

//

type CustomError struct {
	errCode     int
	httpStatus  int
	description string
}

func (c *CustomError) Error() string {
	return c.description
}

func (c *CustomError) ErrorCode() int {
	return c.errCode
}

func (c *CustomError) HttpStatus() int {
	return c.httpStatus
}
