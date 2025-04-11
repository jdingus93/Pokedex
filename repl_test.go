package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "testing 1, 2, 3...",
			expected: []string{"testing", "1,", "2,", "3..."},
		},
		{
			input: "pLEASE wORK",
			expected: []string{"please", "work"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d",
				c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) word %d: got %q, expected %q",
					c.input, i, actual[i], c.expected[i])
			}
		}
	}
}