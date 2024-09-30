package translation

import (
	"errors"
	"testing"
)

func TestString(t *testing.T) {
	testCases := []struct {
		expected   *Varriable
		testString string
		error      error
	}{
		{
			expected: &Varriable{
				Identifier: "123",
				Value:      "test string",
			},
			testString: "@123 ~test string~",
			error:      nil,
		},
		{
			expected:   nil,
			testString: "123 ~test string~",
			error:      ErrIdentifier,
		},
		{
			expected:   nil,
			testString: "@123 test string~",
			error:      ErrValue,
		},
	}
	for _, tc := range testCases {
		start := 0
		v, err := ParseVarriable(tc.testString, &start)
		if !errors.Is(err, tc.error) {
			t.Fatalf("Unexpected Error:\n%+v\n%+v", err, tc.error)
		}
		if tc.expected != nil && *v != *tc.expected {
			t.Fatalf("Failed:\n%+v\n%+v", v, tc.expected)
		}
	}
}
