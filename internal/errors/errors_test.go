package errors

import (
    stdErrors "errors"
    "testing"
)

func TestACBError_Basics(t *testing.T) {
    err := NewACBError(ErrorCodeValidationError, "invalid input").
        WithDetails("field", "name").
        WithError(stdErrors.New("bad value"))

    if err.Code != ErrorCodeValidationError {
        t.Fatalf("expected code %s, got %s", ErrorCodeValidationError, err.Code)
    }
    if err.Details["field"] != "name" {
        t.Fatalf("expected detail field=name, got %v", err.Details)
    }
    if err.Error() == "" {
        t.Fatalf("expected non-empty error string")
    }

    if !Is(err, ErrorCodeValidationError) {
        t.Fatalf("expected Is(..., ValidationError) to be true")
    }
    if Is(err, ErrorCodeUnauthorized) {
        t.Fatalf("expected Is(..., Unauthorized) to be false")
    }
}

func TestConvenienceErrors(t *testing.T) {
    if Unauthorized("u").Code != ErrorCodeUnauthorized {
        t.Fatal("unauthorized code mismatch")
    }
    if Forbidden("f").Code != ErrorCodeForbidden {
        t.Fatal("forbidden code mismatch")
    }
    if NotFound("n").Code != ErrorCodeNotFound {
        t.Fatal("notfound code mismatch")
    }
    if ValidationError("v").Code != ErrorCodeValidationError {
        t.Fatal("validation code mismatch")
    }
    if RateLimitExceeded("r").Code != ErrorCodeRateLimitExceeded {
        t.Fatal("rate limit code mismatch")
    }
    if InternalError("i").Code != ErrorCodeInternalError {
        t.Fatal("internal code mismatch")
    }
    if ServiceUnavailable("s").Code != ErrorCodeServiceUnavailable {
        t.Fatal("service unavailable code mismatch")
    }
}


