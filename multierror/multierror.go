package multierror

import (
	"errors"
	"strings"
	"sync"
)

// MultiError is a collection of errors that implements the error interface.
type MultiError struct {
	mu   sync.RWMutex
	errs []error
}

// NewMultiError returns a new MultiError.
func NewMultiError() *MultiError {
	return &MultiError{}
}

// Error returns the list of errors separated by newlines.
func (e *MultiError) Error() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var errs []string
	for _, err := range e.errs {
		errs = append(errs, err.Error())
	}

	return strings.Join(errs, "\n")
}

// Errors returns the error slice containing the error collection.
func (e *MultiError) Errors() []error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.errs
}

// Add appends an error to the error collection.
func (e *MultiError) Add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Unwrap *MultiError to ensure that depth never exceeds 1.
	var mErr2 *MultiError
	if errors.As(err, &mErr2) {
		e.errs = append(e.errs, mErr2.Errors()...)
		return
	}

	e.errs = append(e.errs, err)
}

// Empty returns whether the *MultiError contains any errors.
func (e *MultiError) Empty() bool {
	return len(e.errs) == 0
}
