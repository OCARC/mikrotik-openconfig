package main

import (
	"testing"
)

// WithDeviceSetting is a helper for safe device state testing.
// It saves the original value, sets a test value, verifies, runs test logic, restores, and verifies restoration.
func WithDeviceSetting(
	t *testing.T,
	getFunc func() (string, error),
	setFunc func(string) error,
	verifyFunc func(string) bool,
	testValue string,
	testLogic func(),
) {
	original, err := getFunc()
	if err != nil {
		t.Fatalf("get original: %v", err)
	}
	if testing.Verbose() {
		t.Logf("ORIGINAL VALUE: %q", original)
	}
	defer func() {
		if err := setFunc(original); err != nil {
			t.Errorf("failed to restore original: %v", err)
		}
		restored, _ := getFunc()
		if testing.Verbose() {
			t.Logf("RESTORED VALUE: %q", restored)
		}
		if !verifyFunc(original) {
			t.Errorf("failed to verify restoration of original value")
		}
	}()
	if err := setFunc(testValue); err != nil {
		t.Fatalf("set test value: %v", err)
	}
	if !verifyFunc(testValue) {
		t.Fatalf("failed to verify test value was set")
	}
	testLogic()
}
