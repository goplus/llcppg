package errs

type AnonymousFuncNotSupportError struct {
}

func (p *AnonymousFuncNotSupportError) Error() string {
	return "anonymous function not supported"
}

func NewAnonymousFuncNotSupportError() *AnonymousFuncNotSupportError {
	return &AnonymousFuncNotSupportError{}
}

type ModNotFoundError struct {
}

func (p *ModNotFoundError) Error() string {
	return "go.mod not found"
}

func NewModNotFoundError() *ModNotFoundError {
	return &ModNotFoundError{}
}
