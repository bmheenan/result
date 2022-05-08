package result_test

import (
	"testing"

	"github.com/bmheenan/result"
)

func TestHandlersToErrReturnError(t *testing.T) {
	err := handlersToErrReturnErr()
	if x, g := "result error: Test error", err.Error(); x != g {
		t.Errorf("Expected error %v; got %v", x, g)
	}
}

func handlersToErrReturnErr() (err error) {
	defer result.HandleError(&err)
	result.Errorf("Test error").
		OrError("result error")
	return nil
}

func TestHandlersToErrReturnNil(t *testing.T) {
	err := handlersToErrReturnNil()
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
}

func handlersToErrReturnNil() (err error) {
	defer result.HandleError(&err)
	result.Ok().OrError("Unexpected error")
	return nil
}

func TestHandlersHandleErrOrReturn(t *testing.T) {
	handlersHandleErrOrDoAndReturn(t)
}

func handlersHandleErrOrDoAndReturn(t *testing.T) (err error) {
	defer result.HandleError(&err)
	result.Errorf("Expected error").
		OrDoAndReturn(func(e error) {})
	t.Error("Coude executed that shouldn't have")
	return nil
}

func TestBasicHandlersHandleReturn(t *testing.T) {
	defer result.HandleReturn()
	result.Errorf("Expected error").
		OrDoAndReturn(func(e error) {})
	t.Error("Code executed that shouldn't have")
}

func TestBasicHandlersHandleReturnUnused(t *testing.T) {
	defer result.HandleReturn()
}
