package analyzer

import "fmt"

type UnreachablePathError struct {
	Path string
	Err  error
}

func (e *UnreachablePathError) Error() string {
	return fmt.Sprintf("Path inaccessible : %s (%v)", e.Path, e.Err)
}

func (e *UnreachablePathError) Unwrap() error {
	return e.Err
}
