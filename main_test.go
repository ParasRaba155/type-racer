package main

import (
	"math"
	"testing"
)

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

func TestTextDiffRatio(t *testing.T) {
	testCases := map[string]struct {
		inputOrig string
		inputNew  string
		output    float64
	}{
		"100% match": {
			inputOrig: "Hello",
			inputNew:  "Hello",
			output:    1.0,
		},
		"0% match": {
			inputOrig: "Hello",
			inputNew:  "asdad",
			output:    0,
		},
		"50% match": {
			inputOrig: "AAAAAA",
			inputNew:  "aAaAaA",
			output:    0.5,
		},
		"99% match": {
			inputOrig: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			inputNew:  "Aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			output:    0.99,
		},
	}

	almostEqual := func(a, b float64) bool {
		const float64EqualityThreshold = 1e-9
		return math.Abs(a-b) <= float64EqualityThreshold
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ratio := textDiffRatio([]byte(test.inputOrig), []byte(test.inputNew))
			if !almostEqual(ratio, test.output) {
				t.Errorf("got %.5f want %.5f", ratio, test.output)
			}
		})
	}
}
