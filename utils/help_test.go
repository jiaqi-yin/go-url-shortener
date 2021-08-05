package utils

import (
	"testing"
)

func TestToSha1(t *testing.T) {
	encodedString := ToSha1("https://www.google.com")
	want := "ef7efc9839c3ee036f023e9635bc3b056d6ee2db"
	if encodedString != want {
		t.Errorf("Got %v, want %v", encodedString, want)
	}
}
