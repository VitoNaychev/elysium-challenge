package assert

import (
	"reflect"
	"testing"
)

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
