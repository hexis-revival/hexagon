package common

import "fmt"

type ErrorCollection struct {
	errors   []error
	position int
}

func (ec *ErrorCollection) Add(err error) {
	if err == nil {
		return
	}
	ec.errors = append(ec.errors, err)
}

func (ec *ErrorCollection) Pop(index int) error {
	if !ec.HasErrors() {
		return nil
	}
	err := ec.errors[index]
	ec.errors = append(ec.errors[:index], ec.errors[index+1:]...)
	return err
}

func (ec *ErrorCollection) Next() error {
	if !ec.HasErrors() {
		return nil
	}
	if ec.position >= len(ec.errors) {
		return nil
	}
	err := ec.errors[ec.position]
	ec.position++
	return err
}

func (ec *ErrorCollection) Errors() []error {
	return ec.errors
}

func (ec *ErrorCollection) Length() int {
	return len(ec.errors)
}

func (ec *ErrorCollection) String() string {
	return fmt.Sprintf("ErrorCollection{errors: %v}", len(ec.errors))
}

func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.errors) > 0
}

func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{}
}

func HandlePanic(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("panic: %v", r)
	}
}
