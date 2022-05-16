package result

import (
	"fmt"
)

// Val is a result that holds one value when ok. Otherwise, it holds an error. It's most useful as a return value for a
// function that either returns a value or an error, e.g:
//     func aPlus1() (res result.Val) {
//         defer result.Handle(&res)
//         a := calcA().
//             OrErr("Couldn't calculate a") // returns an error result.Val
//         return result.NewVal(a + 1)       // returns an ok result.Val
//     }
type Val[T any] struct {
	base
	v T
}

// NewVal returns a new ok Val with the given value v
func NewVal[T any](v T) Val[T] {
	return Val[T]{
		v: v,
	}
}

// ValError returns a new Val with the given error
func ValError[T any](err error) Val[T] {
	v := Val[T]{}
	v.err = err
	return v
}

// ValErrorf returns a new Val with an error made from the given string and arguments. s and args should be the same as
// what would be provided to fmt.Errorf
func ValErrorf[T any](s string, args ...any) Val[T] {
	v := Val[T]{}
	v.err = fmt.Errorf(s, args...)
	return v
}

// TryVal encloses a function that returns a value and an error, then returns its result as a Val. Usage:
//     a := result.TryVal(f()).
//         OrError("f failed")
func TryVal[T any](v T, err error) Val[T] {
	if err == nil {
		return NewVal(v)
	}
	return ValError[T](err)
}

// FromSlice returns a Val containing the value from slice s at position i, if i is within the bounds of s. If i is out
// of bounds, FromSlice returns an error Val
func FromSlice[T any](s []T, i int) Val[T] {
	if i < 0 || i >= len(s) {
		return ValErrorf[T]("Index %v out of bounds for slice of len %v", i, len(s))
	}
	return NewVal(s[i])
}

// FromMap returns a Val containing the value from map m for key k, if there is one. If m has no value for key k,
// FromMap returns an error Val
func FromMap[T any, K comparable](m map[K]T, k K) Val[T] {
	v, ok := m[k]
	if !ok {
		return ValErrorf[T]("Map had no value for key %v", k)
	}
	return NewVal(v)
}

// OrError returns the underlying value if the Val is ok. Otherwise, it stops execution of the calling function and
// returns an error. Use e to provide an explanation about what went wrong; it will be included in the returned error.
//
// OrError must only be used inside a function that returns an error or a result, and that has already defered Handle or
// HandleError. Usage:
//     func aPlus1() (res result.Val[int]) {
//         defer result.Handle(&res)
//         a := calcA().
//             OrError("Couldn't calculate a") // returns an error result.Val
//         return result.NewVal(a + 1)         // returns an ok result.Val
//     }
// If you use OrError without defering Handle or HandleError at the beginning of the function, it will panic
func (v Val[T]) OrError(e string) T {
	if v.err == nil {
		return v.v
	}
	panic(panicToError{
		err: fmt.Errorf("%v: %v", e, v.err),
	})
}

// OrDoAndReturn returns the underlying value if the Val is ok. Otherwise, it executes the provided function f, then
// returns from the calling function.
//
// OrDoAndReturn must only be used inside a function that has already defered Handle, HandleError, or HandleReturn.
// Usage:
//     func main() {
//         defer result.HandleReturn()
//         a := calcA().
//             OrDoAndReturn(func(e error) {
//                 fmt.Printf("Couldn't calculate a: %v\n", e)
//             })
//         fmt.Printf("The value of a is: %v\n", a) // Only executes if calcA returns an ok Val
//     }
// If you use OrDoAndReturn without defering Handle, HandleError, or HandleReturn at the beginning of the function, it
// will panic
func (v Val[T]) OrDoAndReturn(f func(error)) T {
	if v.err == nil {
		return v.v
	}
	f(v.err)
	panic(panicToReturn{
		err: v.err,
	})
}

// OrPanic returns the underlying value if the Val is ok. Otherwise, it panics. This panic will not be caught by Handle,
// HandleError, or HandleReturn. Use p to provide extra context about what went wrong; it will be included in the panic.
// Usage:
//     func main() {
//         a := calcA().
//             OrPanic("Couldn't calculate a")
//         fmt.Printf("The value of a is: %v\n", a)
//     }
func (v Val[T]) OrPanic(p string) T {
	if v.err == nil {
		return v.v
	}
	panic(fmt.Errorf("%v: %v", p, v.err))
}

// OrUse returns the underlying value if the Val is ok. Otherwise, it substitutes in the given value s. Usage:
//     func main() {
//         a := calcA().OrUse(-1)
//         fmt.Printf("The value of a is: %v\n", a) // -1 if calcA returned an error Val
//     }
func (v Val[T]) OrUse(s T) T {
	if v.err == nil {
		return v.v
	}
	return s
}
