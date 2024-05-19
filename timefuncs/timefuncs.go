// Package timefuncs contains generalized functions for working with work time
// and breaks
//
// Yes, this packages was shamelessly copied from
// https://github.com/danverbraganza/meeting-cost-clock
package timefuncs

import (
	"fmt"
	"math"
	"time"
)

func init() {}

func FormatDuration(d time.Duration) string {
	rounder := math.Floor
	if math.Signbit(d.Hours()) {
		rounder = math.Ceil
	}

	return fmt.Sprintf("%02.0f:%02d:%02d",
		rounder(d.Hours()),
		int(math.Abs(d.Minutes()))%60,
		int(math.Abs(d.Seconds()))%60,
		int(math.Abs(float64(d.Nanoseconds()/1e6)))%1000)
}
