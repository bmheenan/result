package result

import (
	"fmt"
)

// Status is the simplest type of result. It's either ok, or it's an error. It's most useful as a return value for a
// function that doesn't need to return any value, but could encounter an error that must be communicated to the calling
// function, e.g:
//     func foo() (res result.Status) {
//         defer result.Handle(&res)
//         doWork().
//             OrErr("Couldn't do work") // returns an error result.Status
//         return result.Ok()            // returns an ok result.Status
//     }
type Status struct {
	base
}

// Ok returns a new ok Status
func Ok() Status {
	return Status{}
}

// Error returns a new Status with the given error
func Error(err error) Status {
	return Status{
		base{
			err: err,
		},
	}
}

// Errorf returns a new Status with an error made from the given string and arguments. s and args should be the same as
// what would be provided to fmt.Errorf
func Errorf(s string, args ...any) Status {
	return Status{
		base{
			err: fmt.Errorf(s, args...),
		},
	}
}

// Try encloses a function that may return an error, then returns its result as a Status. Usage:
//     result.Try(f()).
//         OrError("f failed")
func Try(err error) Status {
	if err == nil {
		return Ok()
	}
	return Error(err)
}

// OrError does nothing if the Status is ok. Otherwise, it stops execution of the calling function and returns an error.
// Use e to provide an explanation about what went wrong; it will be included in the returned error.
//
// OrError must only be used inside a function that returns an error or a result, and that has already defered Handle or
// HandleError. Usage:
//     func f() (res result.Status) {
//         defer result.Handle(&res)
//         doWork().
//             OrError("Couldn't do work")
//         return result.Ok()
//     }
// If you use OrError without defering Handle or HandleError at the beginning of the function, it will panic
func (s Status) OrError(e string) {
	if s.err == nil {
		return
	}
	panic(panicToError{
		err: fmt.Errorf("%v: %v", e, s.err),
	})
}

// OrDoAndReturn does nothing if the Status is ok. Otherwise, it executes the provided function f, then returns from
// the calling function.
//
// OrDoAndReturn must only be used inside a function that has already defered Handle, HandleError, or HandleReturn.
// Usage:
//     func main() {
//         defer result.HandleReturn()
//         doWork().
//             OrDoAndReturn(func(e error) {
//                 fmt.Printf("Couldn't do work: %v\n", e)
//             })
//         fmt.Println("This line only executes if doWork's Status is ok")
//     }
// If you use OrDoAndReturn without defering Handle, HandleError, or HandleReturn at the beginning of the function, it
// will panic
func (s Status) OrDoAndReturn(f func(error)) {
	if s.err == nil {
		return
	}
	f(s.err)
	panic(panicToReturn{
		err: s.err,
	})
}

// OrPanic does nothing if the Status is ok. Otherwise, it panics. This panic will not be caught by Handle,
// HandleError, or HandleReturn. Use p to provide extra context about what went wrong; it will be included in the panic.
// Usage:
//     func main() {
//         doWork().
//             OrPanic("Couldn't do work")
//     }
func (s Status) OrPanic(p string) {
	if s.err == nil {
		return
	}
	panic(fmt.Errorf("%v: %v", p, s.err))
}

// OrDo does nothing if the Status is ok. Otherwise, it executes the provided function f. Usage:
//     func main() {
//         doWork().OrDo(func(e error) {
//             fmt.Printf("Couldn't do work: %v\n", e)
//         })
//     }
func (s Status) OrDo(f func(error)) {
	if s.err == nil {
		return
	}
	f(s.err)
}
