package errors

import "errors"

var (
	// ErrInvalidCEP indica que o CEP fornecido é inválido
	ErrInvalidCEP = errors.New("invalid zipcode")

	// ErrCannotFindZipCode indica que o CEP não foi encontrado
	ErrCannotFindZipCode = errors.New("cannot find zipcode")
)
