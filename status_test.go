package result_test

import (
	"errors"
	"testing"

	"github.com/bmheenan/result"
	"github.com/stretchr/testify/assert"
)

func TestStatusOk(t *testing.T) {
	statusOk().OrDo(func(e error) {
		t.Errorf("statusOk() returned an error status: %v", e)
	})
}

func TestStatusErr(t *testing.T) {
	v := statusErr()
	if v.Ok() {
		t.Error("statusErr() didn't return an error")
	}
}

func TestStatusErrorf(t *testing.T) {
	v := statusErrorf()
	if x, g := "hello world! 1", v.Error(); x != g {
		t.Errorf("Expected error %v; got %v", x, g)
	}
}

func statusErrorf() (v result.Status) {
	defer result.Handle(&v)
	return result.Errorf("%v %v! %d", "hello", "world", 1)
}

func TestStatusOrErrWithError(t *testing.T) {
	v := statusOrErrWithError()
	if x, g := "Error from statusErr: Test error", v.Error(); x != g {
		t.Errorf("Expected error '%v'; got '%v'", x, g)
	}
}

func statusOrErrWithError() (v result.Status) {
	defer result.Handle(&v)
	statusErr().
		OrError("Error from statusErr")
	return result.Ok()
}

func TestStatusOrErrWithoutError(t *testing.T) {
	statusOrErrWithoutError().OrDo(func(e error) {
		t.Error("Error from statusOrErrWithoutError")
	})
}

func statusOrErrWithoutError() (v result.Status) {
	defer result.Handle(&v)
	statusOk().
		OrError("Error from statusOk")
	return result.Ok()
}

func TestStatusOrReturnWithError(t *testing.T) {
	statusOrReturnWithError().OrDo(func(e error) {
		t.Errorf("statusOrReturnWithError returned an error: %v", e)
	})
}

func statusOrReturnWithError() (v result.Status) {
	defer result.Handle(&v)
	statusErr().
		OrDoAndReturn(func(err error) {
			v = result.Ok()
		})
	return result.Errorf("Code executed that should be unreachable")
}

func TestStatusOrReturnWithoutError(t *testing.T) {
	statusOrReturnWithoutError().OrDo(func(e error) {
		t.Errorf("statusOrReturnWithoutError returned an error: %v", e)
	})
}

func statusOrReturnWithoutError() (v result.Status) {
	defer result.Handle(&v)
	statusOk().
		OrDoAndReturn(func(e error) {
			v = result.Errorf("OrReturn executed when it shouldn't have")
		})
	return result.Ok()
}

func TestStatusPanicPassthrough(t *testing.T) {
	defer func() {
		r := recover()
		if r != "Expected panic" {
			t.Errorf("Panic wasn't the expected one, it was: %v", r)
		}
	}()
	statusPanicPassthrough()
	t.Error("Code executed that should be unreachable")
}

func statusPanicPassthrough() (v result.Status) {
	defer result.Handle(&v)
	panic("Expected panic")
}

func TestStatusOrPanicWithError(t *testing.T) {
	defer func() {
		r := recover()
		e, ok := r.(error)
		if !ok {
			t.Errorf("panic wasn't an error, it was a %T", e)
		}
		if x, g := "Expected panic: Test error", e.Error(); x != g {
			t.Errorf("Panic wasn't the expected one, it was: %v", r)
		}
	}()
	statusOrPanicWithError()
	t.Error("Code executed that should be unreachable")
}

func statusOrPanicWithError() (v result.Status) {
	defer result.Handle(&v)
	statusErr().
		OrPanic("Expected panic")
	return result.Ok()
}

func TestStatusOrPanicWithoutError(t *testing.T) {
	statusOk().
		OrPanic("Unexpected panic")
}

func TestStatusOr(t *testing.T) {
	statusOrDo().OrDo(func(e error) {
		t.Errorf("statusOr returned error: %v", e)
	})
}

func statusOrDo() (v result.Status) {
	defer result.Handle(&v)
	v = result.Errorf("This error should be overwritten with an Ok")
	statusErr().OrDo(func(e error) {
		v = result.Ok()
	})
	return
}

func statusOk() (v result.Status) {
	defer result.Handle(&v)
	return result.Ok()
}

func statusErr() (v result.Status) {
	defer result.Handle(&v)
	return result.Error(errors.New("Test error"))
}

func TestTryOk(t *testing.T) {
	result.Try(nil).OrDo(func(e error) {
		t.Errorf("OrDo executed when it shouldn't have")
	})
}

func TestTryError(t *testing.T) {
	a := true
	result.Try(errors.New("Test error")).OrDo(func(e error) {
		a = false
	})
	if a {
		t.Error("OrDo wasn't executed on a Try from an error")
	}
}

func TestErrorIsEmpty(t *testing.T) {
	assert.Equal(t, "", result.Ok().Error())
}
