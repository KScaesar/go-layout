package utility

func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		customCodeChecker: make(map[int]bool),
		errors:            make(map[int]*CustomError),
	}
}

type ErrorRegistry struct {
	customCodeChecker map[int]bool
	errors            map[int]*CustomError
	target            *CustomError
}

func (r *ErrorRegistry) Register(customCode int) *ErrorRegistry {
	if r.target != nil {
		panic("The error registry call sequence is incorrect. It must start with Register and end with Error.")
	}
	if r.customCodeChecker[customCode] {
		panic("duplicated custom code")
	}

	r.customCodeChecker[customCode] = true
	err := &CustomError{customCode: customCode}
	r.errors[err.customCode] = err
	r.target = err
	return r
}

func (r *ErrorRegistry) HttpStatus(httpStatus int) *ErrorRegistry {
	if r.target == nil {
		panic("The error registry call sequence is incorrect. It must start with Register and end with Error.")
	}
	r.target.httpStatus = httpStatus
	return r
}

func (r *ErrorRegistry) Description(description string) *ErrorRegistry {
	if r.target == nil {
		panic("The error registry call sequence is incorrect. It must start with Register and end with Error.")
	}
	r.target.description = description
	return r
}

func (r *ErrorRegistry) Error() error {
	if r.target == nil {
		panic("The error registry call sequence is incorrect. It must start with Register and end with Error.")
	}
	err := r.target
	r.target = nil
	return err
}

//

type CustomError struct {
	customCode  int
	httpStatus  int
	description string
}

func (b *CustomError) Error() string {
	return b.description
}

func (b *CustomError) CustomCode() int {
	return b.customCode
}

func (b *CustomError) HttpStatus() int {
	return b.httpStatus
}
