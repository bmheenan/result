package result_test

import (
	"errors"
	"testing"

	"github.com/bmheenan/result"
	"github.com/stretchr/testify/assert"
)

func TestNewIntVal(t *testing.T) {
	i := result.NewVal(0).
		OrPanic("Couldn't get i")
	assert.Equal(t, 0, i)
}

func TestNewStringVal(t *testing.T) {
	s := result.NewVal("hello").
		OrPanic("Couldn't get s")
	assert.Equal(t, "hello", s)
}

func TestValError(t *testing.T) {
	defer result.HandleReturn()
	_ = result.ValError[int](errors.New("Expected error")).
		OrDoAndReturn(func(e error) {
			assert.Equal(t, "Expected error", e.Error())
		})
	t.Error("This line should not execute")
}

func TestValErrorf(t *testing.T) {
	defer result.HandleReturn()
	_ = result.ValErrorf[int]("Test error %v %v", "hello world", 1).
		OrDoAndReturn(func(e error) {
			assert.Equal(t, "Test error hello world 1", e.Error())
		})
	t.Error("This line should not execute")
}

func TestTryVal(t *testing.T) {
	a := result.TryVal(tryReturnsNil()).OrUse("default")
	assert.Equal(t, "from func", a)
	b := result.TryVal(tryReturnsErr()).OrUse("default")
	assert.Equal(t, "default", b)
}

func tryReturnsNil() (string, error) {
	return "from func", nil
}

func tryReturnsErr() (string, error) {
	return "from func", errors.New("Expected error")
}

func TestValOrErrorResult(t *testing.T) {
	defer result.HandleReturn()
	valOrErrorResult().OrDoAndReturn(func(e error) {
		assert.Equal(t, "Context: Expected error", e.Error())
	})
	t.Error("This not should not execute")
}

func valOrErrorResult() (res result.Val[int]) {
	defer result.Handle(&res)
	result.ValErrorf[int]("Expected error").
		OrError("Context")
	return result.NewVal(0)
}

func TestValOrErrorOkResult(t *testing.T) {
	result.NewVal("hello").
		OrError("Unexpected error")
}

func TestValOrDoAndReturn(t *testing.T) {
	result.NewVal(map[int]string{}).
		OrDoAndReturn(func(e error) {
			t.Errorf("Unexpected error from NewVal: %v", e)
		})
}

func TestValOrPanicPanics(t *testing.T) {
	assert.Panics(t, func() {
		result.ValErrorf[bool]("Expected error").
			OrPanic("Expected panic")
	})
}

func TestValOrPanicDoesntPanic(t *testing.T) {
	result.NewVal("hello").
		OrPanic("Unexpected panic")
}
