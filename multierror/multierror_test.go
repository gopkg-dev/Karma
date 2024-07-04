package multierror

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMultiError(t *testing.T) {
	mErr := NewMultiError()

	if !mErr.Empty() {
		t.Fatal("Empty(): got false, want true")
	}

	var errs []error
	for i := 0; i < 3; i++ {
		err := fmt.Errorf("%d. error", i)
		errs = append(errs, err)
		mErr.Add(err)
	}

	if mErr.Empty() {
		t.Fatal("Empty(): got true, want false")
	}
	if got, want := mErr.Errors(), errs; !reflect.DeepEqual(got, want) {
		t.Errorf("Errors(): got %v, want %v", got, want)
	}

	want := "0. error\n1. error\n2. error"
	if got := mErr.Error(); got != want {
		t.Errorf("Error(): got %q, want %q", got, want)
	}
}
