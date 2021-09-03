package observability

import (
	"fmt"
	"time"
)

// Timer holds start and end times
type Timer struct {
	t1, t2        time.Time
	ela           time.Duration
	timersEnabled bool
	str           string
	on            bool
}

// Start starts the timer
func (timer *Timer) Start(timing bool, str string) {
	if timing {
		timer.t1 = time.Now()
		timer.str = str
	}
}

// EndAndPrint ends the timer and prints the number of seconds
func (timer *Timer) EndAndPrint(timing bool) {
	if timing {
		timer.t2 = time.Now()
		timer.ela = timer.t2.Sub(timer.t1)
		logString := fmt.Sprintf(timer.str+" completed in %.2f ms", float64(timer.ela.Nanoseconds())/1000000.0)
		logN("Info", fmt.Sprintf("%s", logString), 4)
	}
}

// EndAndPrintStderr ends the timer and prints the number of seconds
func (timer *Timer) EndAndPrintStderr(timing bool) {
	if timing {
		timer.t2 = time.Now()
		timer.ela = timer.t2.Sub(timer.t1)
		logString := fmt.Sprintf(timer.str+" completed in %.2f ms", float64(timer.ela.Nanoseconds())/1000000.0)
		logN("Warn", fmt.Sprintf("%s", logString), 4)
	}
}
