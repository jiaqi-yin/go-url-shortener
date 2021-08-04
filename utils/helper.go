package utils

import (
	"crypto/sha1"
	"fmt"
)

func ToSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", string(h.Sum(nil)))
}
