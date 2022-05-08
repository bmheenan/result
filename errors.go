package result

type panicToReturn struct {
	err error
}

func (p panicToReturn) Error() string {
	return "Unrecovered panic from result. Use `defer result.HandleReturn()`, `defer result.Handle(&r)`, or `defer result.HandleError(&err)` at the top of the func to convert the panic into a return"
}

type panicToError struct {
	err error
}

func (p panicToError) Error() string {
	return "Unrecovered panic from result. Use `defer result.Handle(&r)` or `defer result.HandleError(&err)` at the top of the func to convert the panic into a returned result or error: " + p.err.Error()
}
