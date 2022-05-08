package result

import (
	"fmt"
)

// Status is the simplest type of result. It's either ok, or it's an error. It's most useful as a return value for a
// function that doesn't need to return any value, but could encounter an error that must be communicated to the calling
// function, e.g:
//     func foo() (v result.Status) {
//	       defer result.HandleStatus(&v)
//         doWork().
//             OrErr("Couldn't do work") // returns an error result.Status
//         return result.Ok()               // returns an ok result.Status
//     }
type Status struct {
	base
}

// Ok returns a new ok Status
func Ok() Status {
	return Status{}
}

// Err returns a new status with the given error
func Error(err error) Status {
	return Status{
		base{
			err: err,
		},
	}
}

// Errorf returns a new status with an error made from the given string and arguments. s and args should be the same as
// what would be provided to fmt.Errorf
func Errorf(s string, args ...any) Status {
	return Status{
		base{
			err: fmt.Errorf(s, args...),
		},
	}
}

// Try encloses a function that may return an error and returns its result as a Status. Example usage:
//     result.Try(f()).OrError("f failed")
func Try(err error) Status {
	if err == nil {
		return Ok()
	}
	return Error(err)
}

// OrErr does nothing if the Status is ok. Otherwise, it stops execution of the function and returns an error. Use e
// to provide extra context about what went wrong; it will be included in the error.
//
// OrErr must only be used inside a function that returns an error or a result, and that has already set up a Handle
// fuction, e.g:
//     func f() (v result.Status) {
//	       defer opt.HandleStatus(&v)
//         result.Errorf("Test error").
//             OrErr("Found an error") // f will stop executing here and return an error
//         // ... code here will not execute
//     }
//
// If OrErr is called on an Opt without a concrete value, and the function hasn't deferred opt.Handle, it will panic
func (s Status) OrError(e string) {
	if s.err == nil {
		return
	}
	panic(panicToError{
		err: fmt.Errorf("%v: %v", e, s.err),
	})
}

// OrDoAndReturn does nothing if Status is ok. Otherwise, it runs the given function f, and then returns from the
// calling function
func (s Status) OrDoAndReturn(f func(error)) {
	if s.err == nil {
		return
	}
	f(s.err)
	panic(panicToReturn{
		err: s.err,
	})
}

// OrPanic does nothing if Status is ok. Otherwise, if Status has an error, it panics. (This panic will not be caught by
// any result.Handle function). Use p to provide extra context about what went wrong; it will be included in the panic
func (s Status) OrPanic(p string) {
	if s.err == nil {
		return
	}
	panic(fmt.Errorf("%v: %v", p, s.err))
}

// Or does nothing if Status was ok. Otherwise, if Status has an error, it executes the given function, which is given
// access to the error.
func (s Status) OrDo(f func(error)) {
	if s.err == nil {
		return
	}
	f(s.err)
}
