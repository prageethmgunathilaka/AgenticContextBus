package acb

import (
	stdErrors "errors"
	"testing"
)

func TestSDKErrorWrapAndUnwrap(t *testing.T) {
	base := stdErrors.New("boom")
	e := NewSDKError("CODE", "message").WithError(base)
	if e.Error() == "" {
		t.Fatal("expected error string")
	}
	var out error = e
	if !stdErrors.Is(out, base) {
		t.Fatal("unwrap failed")
	}
}
