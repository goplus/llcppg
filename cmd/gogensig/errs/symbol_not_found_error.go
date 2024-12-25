package errs

import (
	"fmt"
)

type SymbolNotFoudError struct {
	name string
}

func (p *SymbolNotFoudError) Error() string {
	return fmt.Sprintf("%s symbol not found", p.name)
}

func NewSymbolNotFoudError(name string) *SymbolNotFoudError {
	return &SymbolNotFoudError{name: name}
}

type SymbolTableNotInitializedError struct {
}

func (p *SymbolTableNotInitializedError) Error() string {
	return "symbol table not initialized"
}

func NewSymbolTableNotInitializedError() *SymbolTableNotInitializedError {
	return &SymbolTableNotInitializedError{}
}
