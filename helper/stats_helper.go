package helper

import (
	"sync"
	"sync/atomic"
	"time"
)

// This is currently only for CPM
// it will have more use in the future

var (
	checksCompleted int64
	cpmStartTime    int64
	finalCPMCount   int64
	endCPMCounter   atomic.Bool
	cpmMutex        sync.Mutex
)

func InitializeCPM() {
	cpmStartTime = time.Now().UnixNano()
	atomic.StoreInt64(&checksCompleted, 0)
}

func EndCPMCounter() {
	atomic.StoreInt64(&finalCPMCount, GetCPM())
	endCPMCounter.Store(true)
}

func IncrementCheckCount() {
	atomic.AddInt64(&checksCompleted, 1)
}

// GetCPM calculates and returns the current Checks Per Minute.
func GetCPM() int64 {
	if endCPMCounter.Load() {
		return finalCPMCount
	}

	cpmMutex.Lock()
	defer cpmMutex.Unlock()

	elapsed := time.Now().UnixNano() - cpmStartTime
	elapsedMinutes := float64(elapsed) / float64(time.Minute.Nanoseconds())

	//This will never divide by zero, because it gets called at least 3 seconds later
	currentCPM := float64(atomic.LoadInt64(&checksCompleted)) / elapsedMinutes

	return int64(currentCPM)
}

func GetChecksCompleted() int64 {
	return atomic.LoadInt64(&checksCompleted)
}
