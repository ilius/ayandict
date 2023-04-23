package search_utils

import (
	"log"
	"time"

	common "github.com/ilius/go-dict-commons"
)

func RunWorkers(
	N int,
	workerCount int,
	timeout time.Duration,
	worker func(int, int) []*common.SearchResultLow,
) []*common.SearchResultLow {
	if workerCount < 2 {
		return worker(0, N)
	}
	if N < 2*workerCount {
		return worker(0, N)
	}

	ch := make(chan []*common.SearchResultLow, workerCount)

	sender := func(start int, end int) {
		ch <- worker(start, end)
	}

	step := N / workerCount
	start := 0
	for i := 0; i < workerCount-1; i++ {
		end := start + step
		go sender(start, end)
		start = end
	}
	go sender(start, N)

	results := []*common.SearchResultLow{}
	timeoutCh := time.NewTimer(timeout)
	for i := 0; i < workerCount; i++ {
		select {
		case wRes := <-ch:
			results = append(results, wRes...)
		case <-timeoutCh.C:
			log.Println("Search Timeout")
			return results
		}
	}

	return results
}
