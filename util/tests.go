package util

import (
	"fmt"
	"testing"
)

func AssertNoError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("Unexpected error: \"%s\"", msg)
	}
}

func AssertNoErrorf(t *testing.T, err error, msg string, args ...any) {
	if err != nil {
		t.Fatalf("Unexpected error: \"%s\"", fmt.Sprintf(msg, args...))
	}
}

func AssertError(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatalf("Expected error: \"%s\"", msg)
	}
}

func AssertErrorf(t *testing.T, err error, msg string, args ...any) {
	if err == nil {
		t.Fatalf("Expected error: \"%s\"", fmt.Sprintf(msg, args...))
	}
}

func AssertEqual(t *testing.T, actual any, expected any, msg string) {
	if actual != expected {
		t.Fatalf("Expected \"%v\", Actual \"%v\" : %s", expected, actual, msg)
	}
}

func AssertEqualf(t *testing.T, actual any, expected any, msg string, args ...any) {
	if actual != expected {
		t.Fatalf("Expected \"%v\", Actual \"%v\" : %s", expected, actual, fmt.Sprintf(msg, args...))
	}
}

func AssertNotEqual(t *testing.T, actual any, expected any, msg string) {
	if actual == expected {
		t.Fatalf("Expected \"%v\" to be different from actual value : %s", expected, msg)
	}
}

func AssertNotEqualf(t *testing.T, actual any, expected any, msg string, args ...any) {
	if actual == expected {
		t.Fatalf("Expected \"%v\" to be different from actual value : %s", expected, fmt.Sprintf(msg, args...))
	}
}
