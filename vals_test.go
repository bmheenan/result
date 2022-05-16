package result_test

import (
	"errors"
	"testing"

	"github.com/bmheenan/result"
	"github.com/stretchr/testify/assert"
)

func TestTryVals(t *testing.T) {
	a, b := result.TryVals(func() (string, string, error) {
		return "hello", "world", nil
	}()).OrPanic("Couldn't get strings")
	assert.Equal(t, "hello world", a+" "+b)

	c, d := result.TryVals(func() (string, string, error) {
		return "", "", errors.New("expected error")
	}()).OrUse("hello", "world")
	assert.Equal(t, "hello world", c+" "+d)
}

func TestValsErrorf(t *testing.T) {
	assert.EqualError(
		t,
		result.ValsErrorf[int, int]("Expected error %v %v", 10, "hi"),
		"Expected error 10 hi",
	)
}

func TestValsOrError(t *testing.T) {
	defer result.HandleReturn()

	a, b := func() (res result.Vals[int, int]) {
		v := result.NewVals(1, 2)
		v.
			OrError("Unexpected error")
		return v
	}().OrPanic("Got an error")
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)

	_, _ = func() (res result.Vals[int, int]) {
		defer result.Handle(&res)
		result.ValsErrorf[int, int]("Expected error").
			OrError("OrError triggered")
		return result.NewVals(0, 0)
	}().OrDoAndReturn(func(err error) {
		assert.EqualError(
			t,
			err,
			"OrError triggered: Expected error",
		)
	})
	t.Errorf("This line should not execute")
}

func TestValsOrDoAndReturn(t *testing.T) {
	defer result.HandleReturn()

	a, b := result.NewVals("hello", "world").
		OrDoAndReturn(func(err error) {
			t.Errorf("This line should not execute")
		})
	assert.Equal(t, "hello", a)
	assert.Equal(t, "world", b)

	_, _ = result.ValsErrorf[int, int]("Expected error").
		OrDoAndReturn(func(err error) {
			assert.EqualError(
				t,
				err,
				"Expected error",
			)
		})
	t.Errorf("This line should not execute")
}

func TestValsOrPanic(t *testing.T) {
	a, b := result.NewVals(1.1, "1.2").
		OrPanic("NewVals was an error Vals")
	assert.Equal(t, 1.1, a)
	assert.Equal(t, "1.2", b)

	assert.PanicsWithErrorf(
		t,
		"Panic: Expected error",
		func() {
			_, _ = result.ValsErrorf[int, string]("Expected error").
				OrPanic("Panic")
		},
		"Expected panic from error Vals",
	)
}

func TestValsOrUse(t *testing.T) {
	a, b := result.NewVals(true, -50).
		OrUse(false, 100)
	assert.Equal(t, true, a)
	assert.Equal(t, -50, b)

	c, d := result.ValsErrorf[map[int]string, int]("Expected error").
		OrUse(map[int]string{
			0: "hello",
		}, 100)
	assert.Equal(t, map[int]string{0: "hello"}, c)
	assert.Equal(t, 100, d)
}
