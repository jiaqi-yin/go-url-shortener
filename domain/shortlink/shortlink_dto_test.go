package shortlink

import (
	"testing"

	"github.com/jiaqi-yin/go-url-shortener/utils"
)

func TestValidateSuccess(t *testing.T) {
	req := &ShortenRequest{
		URL: "http://example.com",
	}
	if err := req.Validate(); err != nil {
		t.Errorf("Got %v, want %v", err, nil)
	}
}

func TestValidateError(t *testing.T) {
	req := &ShortenRequest{}
	want := utils.NewBadRequestError("missing url in the request")
	if err := req.Validate(); err != want {
		t.Errorf("Got %v, want %v", err, want)
	}
}
