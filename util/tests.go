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
		t.Fatalf("Expected \"%v\", Actual \"%v\" (%T) : %s", expected, actual, actual, msg)
	}
}

func AssertEqualf(t *testing.T, actual any, expected any, msg string, args ...any) {
	if actual != expected {
		t.Fatalf("Expected \"%v\", Actual \"%v\" (%T) : %s", expected, actual, actual, fmt.Sprintf(msg, args...))
	}
}

func AssertEqualExists(t *testing.T, actual any, actualExists bool, expected any, msg string) {
	if !actualExists {
		t.Fatalf("Expected a value to exist : %s", msg)
	}
	AssertEqual(t, actual, expected, msg)
}

func AssertEqualExistsf(t *testing.T, actual any, actualExists bool, expected any, msg string, args ...any) {
	if !actualExists {
		t.Fatalf("Expected a value to exist : %s", fmt.Sprintf(msg, args...))
	}
	AssertEqualf(t, actual, expected, msg, args...)
}

func AssertEqualExistsFunc(t *testing.T, getValue func() (any, bool), expected any, msg string) {
	actual, actualExists := getValue()
	AssertEqualExists(t, actual, actualExists, expected, msg)
}

func AssertEqualExistsFuncf(t *testing.T, getValue func() (any, bool), expected any, msg string, args ...any) {
	actual, actualExists := getValue()
	AssertEqualExistsf(t, actual, actualExists, expected, msg, args...)
}

func AssertNotEqual(t *testing.T, value any, shouldntMatchThis any, msg string) {
	if value == shouldntMatchThis {
		t.Fatalf("Expected \"%v\" to be different from actual value : %s", shouldntMatchThis, msg)
	}
}

func AssertNotEqualf(t *testing.T, value any, shouldntMatchThis any, msg string, args ...any) {
	if value == shouldntMatchThis {
		t.Fatalf("Expected \"%v\" to be different from actual value : %s", shouldntMatchThis, fmt.Sprintf(msg, args...))
	}
}

func AssertNoPanic(t *testing.T, test func(), msg string) {
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("Test paniced: %s - %v", msg, err)
			}
		}()
		test()
	}()
}

func AssertPanic(t *testing.T, test func(), msg string) {
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatalf("Test did not panic: %s", msg)
			}
		}()
		test()
	}()
}
