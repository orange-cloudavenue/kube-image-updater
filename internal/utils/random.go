package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func RandomInRange(start, end int) int {
	return rand.Intn(end-start+1) + start
}

func RandomSecondInRange(start, end int) time.Duration {
	return time.Duration(RandomInRange(start, end)) * time.Second
}
