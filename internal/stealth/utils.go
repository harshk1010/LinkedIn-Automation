package stealth

import (
	"math/rand"
	"time"
)

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomDelay(minMs, maxMs int) {
	if minMs <= 0 || maxMs <= 0 || minMs > maxMs {
		return
	}

	delay := rand.Intn(maxMs-minMs+1) + minMs
	time.Sleep(time.Duration(delay) * time.Millisecond)
}