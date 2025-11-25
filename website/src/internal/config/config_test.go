package config

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestGetValue(t *testing.T) {
	expectedKey := "hello"
	expectedValue := "world"

	resultKey, resultValue, err := getValue(expectedKey + "=" + expectedValue)

	if resultKey != expectedKey || resultValue != expectedValue || err != nil {
		t.Errorf("expected '%s=%s' and no error, got '%s=%s' error: %s", expectedKey, expectedValue, resultKey, resultValue, err)
	}
}

func TestGetValueWithNewLine(t *testing.T) {
	expectedKey := "hello"
	expectedValue := "world"

	resultKey, resultValue, err := getValue(expectedKey + "=" + expectedValue + "\n")

	if resultKey != expectedKey || resultValue != expectedValue {
		t.Errorf("expected '%s=%s' and no error, got '%s=%s' error: %s", expectedKey, expectedValue, resultKey, resultValue, err)
	}
}

func TestGetValueEmpty(t *testing.T) {
	saveStdr := os.Stderr
	defer func() {
		os.Stderr = saveStdr
	}()

	var buf bytes.Buffer
	log.SetOutput(&buf)

	resultKey, resultValue, err := getValue("")

	if resultKey != "" || resultValue != "" {
		t.Errorf("expected empty strings and no error, got '%s=%s' error: %s", resultKey, resultValue, err)
	}
}

func TestGetValueNoEqual(t *testing.T) {
	resultKey, resultValue, err := getValue("hello")

	if resultKey != "" || resultValue != "" || err != errEqualSignMissing {
		t.Errorf("expected empty strings and error on equal sign, got '%s=%s' error: %s", resultKey, resultValue, err)
	}
}

func TestGetValueEmptyKey(t *testing.T) {
	expectedKey := ""
	expectedValue := "world"

	resultKey, resultValue, err := getValue(expectedKey + "=" + expectedValue)

	if resultKey != "" || resultValue != "" || err != errKeyEmpty {
		t.Errorf("expected empty strings and error on empty key, got '%s=%s' error: %s", resultKey, resultValue, err)
	}
}

func TestGetValueEmptyValue(t *testing.T) {
	expectedKey := "hello"
	expectedValue := ""

	resultKey, resultValue, err := getValue(expectedKey + "=" + expectedValue)

	if resultKey != expectedKey || resultValue != expectedValue || err != nil {
		t.Errorf("expected '%s=%s' and no error, got '%s=%s' error: %s", expectedKey, expectedValue, resultKey, resultValue, err)
	}
}
