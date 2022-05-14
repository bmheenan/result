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

func TestFromSliceInBounds(t *testing.T) {
	s := []string{"hello", "world"}
	assert.Equal(
		t,
		result.FromSlice(s, 0).
			OrPanic("Couldn't get position 0"),
		"hello",
	)
	assert.Equal(
		t,
		result.FromSlice(s, 1).
			OrPanic("Couldn't get position 1"),
		"world",
	)
}

func TestFromSliceOutOfBounds(t *testing.T) {
	s0 := []string{}
	s := []string{"hello", "world"}
	assert.PanicsWithErrorf(
		t,
		"No value: Index 0 out of bounds for slice of len 0",
		func() {
			result.FromSlice(s0, 0).OrPanic("No value")
		},
		"Didn't panic from len 0",
	)
	assert.PanicsWithErrorf(
		t,
		"No value: Index 2 out of bounds for slice of len 2",
		func() {
			result.FromSlice(s, 2).OrPanic("No value")
		},
		"Didn't panic from len 2",
	)
	assert.PanicsWithErrorf(
		t,
		"No value: Index -1 out of bounds for slice of len 2",
		func() {
			result.FromSlice(s, -1).OrPanic("No value")
		},
		"Didn't panic from negative i",
	)
}

func TestFromMapPresentKey(t *testing.T) {
	m := map[int]string{
		1:   "hello",
		100: "world",
	}
	assert.Equal(
		t,
		result.FromMap(m, 1).
			OrPanic("No value"),
		"hello",
	)
}

func TestFromMapMissingKey(t *testing.T) {
	assert.PanicsWithErrorf(
		t,
		"No value: Map had no value for key 0",
		func() {
			result.FromMap(map[int]string{}, 0).OrPanic("No value")
		},
		"Didn't panic from empty map",
	)
}
