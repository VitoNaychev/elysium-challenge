package assert

import (
	"errors"
	"reflect"
	"testing"
)

func ErrorType[T error](t testing.TB, got error) {
	var want T

	if !errors.As(got, &want) {
		t.Errorf("got error with type %v want %v", reflect.TypeOf(got), reflect.TypeOf(want))
	}
}

func RequireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got error: %v ", err)
	}
}

func Equal[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
