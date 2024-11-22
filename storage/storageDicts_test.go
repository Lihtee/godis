package storage

import (
	"reflect"
	"testing"
	"time"
)

var (
	normalMap     = map[string]string{"a": "1"}
	zeroLengthMap = map[string]string{}
	halfEmptyMap  = map[string]string{"a": "1", "b": NullValue}
	allEmptyMap   = map[string]string{"": ""}
)

func TestSetGetDict(t *testing.T) {
	cases := []struct {
		key            string
		setValue       map[string]string
		setTtl         time.Duration
		getValue       []string
		expectedOutput map[string]string
		expectedError  bool
	}{
		{"spi", normalMap, 0, []string{"a"}, normalMap, false},
		{"spi", zeroLengthMap, 0, []string{}, zeroLengthMap, false},
		{"spi", nil, 0, []string{}, nil, false},
		{"spi", nil, 0, []string{"a"}, nil, false},
		{"spi", allEmptyMap, 0, []string{}, allEmptyMap, false},
		{"spi", allEmptyMap, 0, []string{""}, allEmptyMap, false},
		{"spi", normalMap, 0, []string{"a", "b"}, halfEmptyMap, false},
		{"spi", normalMap, 0, nil, normalMap, false},
	}

	for _, testCase := range cases {
		storage := New(true)
		actual, err := storage.GetDict(testCase.key, testCase.getValue...)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on get: %v", err)
			} else {
				continue
			}
		}
		if actual != nil {
			t.Errorf("expecting nil but got %v", actual)
		}

		err = storage.SetDict(testCase.key, testCase.setValue, testCase.setTtl)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on set: %v", err)
			} else {
				continue
			}
		}

		actual, err = storage.GetDict(testCase.key, testCase.getValue...)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on get: %v", err)
			} else {
				continue
			}
		}
		if !reflect.DeepEqual(actual, testCase.expectedOutput) {
			t.Errorf("expecting %v but got %v", actual, testCase.expectedOutput)
		}

		err = storage.DeleteKey(testCase.key)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on delete: %v", err)
			} else {
				continue
			}
		}

		actual, err = storage.GetDict(testCase.key, testCase.getValue...)
		if err != nil {
			if !testCase.expectedError {
				t.Errorf("unexpected error on get: %v", err)
			} else {
				continue
			}
		}
		if actual != nil {
			t.Errorf("expecting nil but got %v", actual)
		}
	}
}
