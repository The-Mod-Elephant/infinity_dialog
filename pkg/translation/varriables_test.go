package translation

import (
	"errors"
	"testing"
)

func TestString(t *testing.T) {
	testCases := []struct {
		expected   *Variable
		testString string
		error      error
	}{
		{
			expected: &Variable{
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
		{
			expected: &Variable{
				Identifier: "0",
				Value:      "Aboleth le�� ve sv� n�dr�i, mimo dosah tv�ch zbran�. Mus� naj�t jin� zp�sob, jak zni�it mechanismus, kter� ho dr�� na�ivu.",
			},
			testString: "@0    = ~Aboleth le�� ve sv� n�dr�i, mimo dosah tv�ch zbran�. Mus� naj�t jin� zp�sob, jak zni�it mechanismus, kter� ho dr�� na�ivu.~",
			error:      nil,
		},
	}
	for _, tc := range testCases {
		v, err := FromString(tc.testString, 0)
		if !errors.Is(err, tc.error) {
			t.Fatalf("Unexpected Error:\n%+v\n%+v", err, tc.error)
		}
		if tc.expected != nil &&
			v.Identifier != tc.expected.Identifier &&
			v.Value != tc.expected.Value {
			t.Fatalf("Failed:\n%+v\n%+v", v, tc.expected)
		}
	}
}
