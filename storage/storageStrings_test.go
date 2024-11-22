package storage

import (
	"math"
	"testing"
	"time"
)

func TestSetGetString(t *testing.T) {
	cases := []struct {
		inputKey       string
		inputValue     string
		inputTtl       time.Duration
		expectedOutput string
		expectedError  bool
	}{
		{"spi", "1", 0, "1", false},
		{"spi", "", 0, "", false},
		{"", "", 0, "", false},
		{"", "", -1, NullValue, false},
		{"", "", math.MaxInt64, "", false},
		{"", "", math.MinInt64, NullValue, false},
	}

	for _, testCase := range cases {
		storage := New(true)
		err := storage.SetString(testCase.inputKey, testCase.inputValue, testCase.inputTtl)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on set: %v", err)
			} else {
				continue
			}
		}

		actual, err := storage.GetString(testCase.inputKey)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on get: %v", err)
			} else {
				continue
			}
		}

		if actual != testCase.expectedOutput {
			t.Errorf("expected %s, got %s", testCase.expectedOutput, actual)
		}

		err = storage.DeleteKey(testCase.inputKey)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on remove: %v", err)
			} else {
				continue
			}
		}

		actual, err = storage.GetString(testCase.inputKey)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on get after remove: %v", err)
			} else {
				continue
			}
		}

		if actual != NullValue {
			t.Errorf("expected %s, got %s", NullValue, actual)
		}
	}
}
