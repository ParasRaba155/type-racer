package main

import "testing"

func TestWrappedText(t *testing.T) {
	testCases := map[string]struct {
		input string
		width int
		ouput string
	}{
		"Should not wrap": {
			input: "Hello World",
			ouput: "Hello World",
			width: 11,
		},
		"Should be wrapped with a single line": {
			input: "Hello World",
			ouput: "Hello \nWorld",
			width: 6,
		},
		"Should be wrapped with a double line": {
			input: "Hello World",
			ouput: "Hell\no Wo\nrld",
			width: 4,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			got := getWrappedText(test.input, test.width)
			if test.ouput != got {
				t.Errorf("got %q want %q", got, test.ouput)
			}
		})
	}
}
