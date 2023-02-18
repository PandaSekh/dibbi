package dibbi

import "errors"

var (
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrColumnDoesNotExist = errors.New("column does not exist")
	ErrInvalidDatatype    = errors.New("invalid datatype")
	ErrMissingValues      = errors.New("missing values")
	ErrInvalidCell        = errors.New("cell content is invalid")
	ErrInvalidOperands    = errors.New("operands are invalid")
)
