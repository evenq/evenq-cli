package util

import (
	"time"

	"github.com/briandowns/spinner"
)

func Spinner(prefix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[42], 100*time.Millisecond)
	s.Prefix = prefix + " "
	s.Start()

	return s
}
