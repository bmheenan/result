package result_test

import (
	"testing"

	"github.com/bmheenan/result"
)

func TestErrorsOrReturnPanics(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok {
			t.Errorf("Panic wasn't an error, it was a %T", r)
		}
		if x, g := "Unrecovered panic from result. Use `defer result.HandleReturn()`, `defer result.HandleStatus(&v)`, or `defer result.HandleErr(&err)` at the top of the func to convert the panic into a return",
			e.Error(); x != g {

			t.Errorf("Expected Error '%v' but got '%v'", x, g)
		}
	}()
	errorsOrReturnPanics()
}

func errorsOrReturnPanics() {
	errorsErr().
		OrDoAndReturn(func(e error) {})
}

func TestErrorsOrErrPanics(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok {
			t.Errorf("Panic wasn't an error, it was a %T", r)
		}
		if x, g := "Unrecovered panic from result. Use `defer result.HandleStatus(&v)` or `defer result.HandleErr(&err)` at the top of the func to convert the panic into a returned result or error: Expected error: Test error",
			e.Error(); x != g {

			t.Errorf("Expected Error '%v' but got '%v'", x, g)
		}
	}()
	errorsOrErrPanics()
}

func errorsOrErrPanics() {
	errorsErr().
		OrError("Expected error")
}

func errorsErr() (v result.Status) {
	return result.Errorf("Test error")
}
