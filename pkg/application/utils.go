package application

import (
	"strconv"
	"strings"
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

func joinIntList(nums []int) string {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(int64(num), 10)
	}
	return strings.Join(strs, ",")
}

func splitIntList(st string) ([]int, error) {
	strs := strings.Split(st, ",")
	nums := make([]int, len(strs))
	for i, st := range strs {
		n, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			return nil, err
		}
		nums[i] = int(n)
	}
	return nums, nil
}
