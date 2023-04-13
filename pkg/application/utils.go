package application

import (
	"sync"
	"time"
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func tryLockAsManyAs(mutex *sync.Mutex, count int, sleep time.Duration) bool {
	for i := 0; i < count; i++ {
		if mutex.TryLock() {
			return true
		}
		time.Sleep(sleep)
	}
	return false
}
