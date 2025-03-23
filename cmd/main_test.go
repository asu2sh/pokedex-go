package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Hola amigo    !",
			expected: []string{"hola", "amigo", "!"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(c.expected) != len(actual) {
			t.Errorf("Length doesn't match!")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Words don't match!")
			}
		}
	}
}
