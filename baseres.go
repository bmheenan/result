package result

type errorSetter interface {
	setError(error)
}

// base holds the basic functionality shared by all results
type base struct {
	err error
}

// Ok returns whether the result is ok. If not, it will have an error
func (b base) Ok() bool {
	return b.err == nil
}

func (b *base) setError(err error) {
	b.err = err
}

// Err returns the error if the result has one. Otherwise, if the result is ok, it returns ""
func (b base) Error() string {
	if b.err == nil {
		return ""
	}
	return b.err.Error()
}
