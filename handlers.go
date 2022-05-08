package result

// HandleReturn must be defered at the beginning of a function if that function doesn't return an error or a result, in
// order to use .OrReturn within the function
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

// HandleErr must be used with defer at the start of a function if that function returns an error. If the function
// uses .OrErr or .OrReturn without defering a HandleXxx function, it will panic. err must be a pointer to the named
// error return value of the function. Usage:
//     func f() (err error) {
//	       defer result.HandleErr(&err)
//         // now safe to use:
//         x := result.None[int]().OrErr("Couldn't get x") // f will stop executing here and return an error
//         // ...
//     }
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

// Handle must be deferred at the beginning of a function, if that function returns a result, in order to use
// .OrReturn or .OrError within the function
func Handle(e errorSetter) {
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
		e.setError(p.err)
		return
	}
	panic(r)
}
