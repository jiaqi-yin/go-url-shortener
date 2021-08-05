package utils

import (
	"net/http"
	"testing"
)

func TestRestErr(t *testing.T) {
	err := &restErr{
		ErrMessage: "invalid json",
		ErrStatus:  400,
		ErrError:   "bad_request",
	}
	if msg := err.Message(); msg != "invalid json" {
		t.Errorf("Got %v, want %v", msg, "invalid json")
	}
	if status := err.Status(); status != 400 {
		t.Errorf("Got %v, want %v", status, 400)
	}
	if error := err.Error(); error != "message: invalid json - status: 400 - error: bad_request" {
		t.Errorf("Got %v, want %v", error, "message: invalid json - status: 400 - error: bad_request")
	}
}

func TestNewRestError(t *testing.T) {
	message := "foo"
	status := 400
	err := "bar"
	restErr := NewRestError(message, status, err)
	if !isRestErr(restErr) {
		t.Errorf("Got %T, want %v", restErr, "RestErr")
	}
}

func NewError(f func(string) RestErr, message string) RestErr {
	return f(message)
}

func TestNewError(t *testing.T) {
	tests := []struct {
		f       func(string) RestErr
		message string
		want    RestErr
	}{
		{
			f:       NewBadRequestError,
			message: "bad request error",
			want: restErr{
				ErrMessage: "bad request error",
				ErrStatus:  http.StatusBadRequest,
				ErrError:   "bad_request",
			},
		},
		{
			f:       NewNotFoundError,
			message: "resource not found",
			want: restErr{
				ErrMessage: "resource not found",
				ErrStatus:  http.StatusNotFound,
				ErrError:   "not_found",
			},
		},
		{
			f:       NewInternalServerError,
			message: "internal server error",
			want: restErr{
				ErrMessage: "internal server error",
				ErrStatus:  http.StatusInternalServerError,
				ErrError:   "internal_server_error",
			},
		},
	}
	for _, tt := range tests {
		if got := tt.f(tt.message); got != tt.want {
			t.Errorf("Got %v, want %v", got, tt.want)
		}
	}
}

func isRestErr(t interface{}) bool {
	switch t.(type) {
	case RestErr:
		return true
	default:
		return false
	}
}
