package utils

import "time"

func DoEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}
