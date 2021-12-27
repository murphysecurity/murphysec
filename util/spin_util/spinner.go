//go:build !noprint

package spin_util

import (
	"github.com/briandowns/spinner"
	"time"
)

var s = func() *spinner.Spinner { return spinner.New(spinner.CharSets[43], 100*time.Millisecond) }()

func StartSpinner(prefix string, suffix string) {
	s.Stop()
	s.Prefix = prefix
	s.Suffix = suffix
	s.Start()
}

func StopSpinner() {
	s.Stop()
}
