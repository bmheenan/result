package result

// HandleReturn must be defered at the beginning of a function if that function doesn't return an error or a result, in
// order to use OrDoAndReturn within the function. Usage:
//     func main() {
//         defer result.HandleReturn()
//         // now safe to use:
//         result.Errorf("Result with error").
//             OrDoAndReturn(func(e error) {
//                 fmt.Println("Result contained an error") // main will stop executing after this line
//             })
//         fmt.Println("this line will not execute")
//     }
// If you use OrError or OrDoAndReturn without defering Handle, HandleError, or HandleReturn at the beginning of the
// function, it will panic
func HandleReturn() {
	r := recover()
	if r == nil {
		return
	}
	_, ok := r.(panicToReturn)
	if ok {
		return
	}
	panic(r)
}

// HandleError must be defered at the begining of a function if that function returns an error, in order to to use
// OrError or OrDoAndReturn within the function. err must be a pointer to the named error return value of the
// function. Usage:
//     func f() (err error) {
//         defer result.HandleError(&err)
//         // now safe to use:
//         result.Errorf("Result with an error").
//             OrError("Result contained an error") // f will stop executing here and return an error
//         fmt.Println("this line will not execute")
//     }
// If you use OrError or OrDoAndReturn without defering Handle, HandleError, or HandleReturn at the beginning of the
// function, it will panic
func HandleError(err *error) {
	r := recover()
	if r == nil {
		return
	}
	_, ok := r.(panicToReturn)
	if ok {
		return
	}
	e, ok := r.(panicToError)
	if ok {
		*err = e.err
		return
	}
	panic(r)
}

// Handle must be defered at the begining of a function if that function returns a result, in order to to use
// OrError or OrDoAndReturn within the function. res must be a pointer to the named result return value of the
// function. Usage:
//     func f() (res result.Status) {
//         defer result.Handle(&res)
//         // now safe to use:
//         result.Errorf("Result with an error").
//             OrError("Result contained an error") // f will stop executing here and return an error Status
//         fmt.Println("this line will not execute")
//     }
// If you use OrError or OrDoAndReturn without defering Handle, HandleError, or HandleReturn at the beginning of the
// function, it will panic
func Handle(res errorSetter) {
	r := recover()
	if r == nil {
		return
	}
	_, ok := r.(panicToReturn)
	if ok {
		return
	}
	p, ok := r.(panicToError)
	if ok {
		res.setError(p.err)
		return
	}
	panic(r)
}
