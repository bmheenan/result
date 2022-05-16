package result

import (
	"fmt"
)

// Vals is a result that holds 2 values when ok. Otherwise, it holds an error. It's most useful as a return value for a
// function that either returns 2 values or an error, e.g:
//     func divAndMod(a, b int) result.Vals {
//         if b == 0 {
//             return result.ValsErrorf("Cannot divide by zero")
//         }
//         return result.NewVals(a / b, a % b)
//     }
type Vals[T, U any] struct {
	base
	v0 T
	v1 U
}

// NewVal returns a new ok Vals with the given values v0, and v1
func NewVals[T, U any](v0 T, v1 U) Vals[T, U] {
	return Vals[T, U]{
		v0: v0,
		v1: v1,
	}
}

// ValsError returns a new Vals with the given error
func ValsError[T, U any](err error) Vals[T, U] {
	v := Vals[T, U]{}
	v.err = err
	return v
}

// ValsErrorf returns a new Vals with an error made from the given string and arguments. s and args should be the same
// as what would be provided to fmt.Errorf
func ValsErrorf[T, U any](s string, args ...any) Vals[T, U] {
	v := Vals[T, U]{}
	v.err = fmt.Errorf(s, args...)
	return v
}

// TryVals encloses a function that returns two values and an error, then returns its result as a Vals. Usage:
//     a, b := result.TryVal(f()).
//         OrError("f failed")
func TryVals[T, U any](v0 T, v1 U, err error) Vals[T, U] {
	if err == nil {
		return NewVals(v0, v1)
	}
	return ValsError[T, U](err)
}

// OrError returns the underlying values if the Vals is ok. Otherwise, it stops execution of the calling function and
// returns an error. Use e to provide an explanation about what went wrong; it will be included in the returned error.
//
// OrError must only be used inside a function that returns an error or a result, and that has already defered Handle or
// HandleError. Usage:
//     func employeeFullName(id int) (res result.Val[int]) {
//         defer result.Handle(&res)
//         first, last := employeeNames(id). // employeeNames returns a result.Vals
//             OrError("Couldn't lookup employee names")
//         return result.NewVal(first + " " + last)
//     }
// If you use OrError without defering Handle or HandleError at the beginning of the function, it will panic
func (v Vals[T, U]) OrError(e string) (T, U) {
	if v.err == nil {
		return v.v0, v.v1
	}
	panic(panicToError{
		err: fmt.Errorf("%v: %v", e, v.err),
	})
}

// OrDoAndReturn returns the underlying values if the Vals is ok. Otherwise, it executes the provided function f, then
// returns from the calling function.
//
// OrDoAndReturn must only be used inside a function that has already defered Handle, HandleError, or HandleReturn.
// Usage:
//     func main() {
//         defer result.HandleReturn()
//         user, pass := parseFlags(). // parseFlags returns a Vals
//             OrDoAndReturn(func(e error) {
//                 fmt.Printf("Couldn't parse flags: %v\n", e)
//             })
//         setup(user, pass)
//     }
// If you use OrDoAndReturn without defering Handle, HandleError, or HandleReturn at the beginning of the function, it
// will panic
func (v Vals[T, U]) OrDoAndReturn(f func(error)) (T, U) {
	if v.err == nil {
		return v.v0, v.v1
	}
	f(v.err)
	panic(panicToReturn{
		err: v.err,
	})
}

// OrPanic returns the underlying values if the Vals is ok. Otherwise, it panics. This panic will not be caught by
// Handle, HandleError, or HandleReturn. Use p to provide extra context about what went wrong; it will be included in
// the panic. Usage:
//     func main() {
//         user, pass := parseFlags(). // parseFlags returns a Vals
//             OrPanic("Couldn't parse flags")
//         setup(user, pass)
//     }
func (v Vals[T, U]) OrPanic(p string) (T, U) {
	if v.err == nil {
		return v.v0, v.v1
	}
	panic(fmt.Errorf("%v: %v", p, v.err))
}

// OrUse returns the underlying values if the Vals is ok. Otherwise, it substitutes in the given values s0 and s1.
// Usage:
//     func main() {
//         user, pass := parseFlags(). // parseFlags returns a Vals
//             OrUse("defaultadmin", "defaultpass_123")
//         setup(user, pass)
//     }
func (v Vals[T, U]) OrUse(s0 T, s1 U) (T, U) {
	if v.err == nil {
		return v.v0, v.v1
	}
	return s0, s1
}
