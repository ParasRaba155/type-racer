package main

import "fmt"

type Stats struct {
	Accuracy float64
	// Word Per Minute
	WPM float64
}

func (s Stats) String() string {
	return fmt.Sprintf("Accuracy = %.0f%%, WPM (Words Per Minute) = %.f", s.Accuracy, s.WPM)
}

// GetStats orig and user input text and total time in second
func GetStats(orig, newer []rune, sec float64) Stats {
	diffRatio := textDiffRatio(orig, newer)

	return Stats{
		Accuracy: diffRatio * 100,
		WPM:      getWPM(newer, sec),
	}
}

// getWPM calculates by considering 5 runes as a single word
func getWPM(str []rune, sec float64) float64 {
	numOfLetterInWord := 5.0
	numOfLetters := 0
	for _, char := range str {
		if char == 0 {
			continue
		}
		numOfLetters++
	}
	numOfWordInStr := float64(numOfLetters) / numOfLetterInWord

	wps := numOfWordInStr / sec
	return wps * 60
}
