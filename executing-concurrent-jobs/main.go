package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	jobs := []Job{
		SleepJob{100 * time.Millisecond, 100},
		SleepJob{300 * time.Millisecond, 300},
		SleepJob{500 * time.Millisecond, 500},
		SleepJob{700 * time.Millisecond, 700},
	}
	fmt.Println(ExecuteWithTimeout(jobs, 500*time.Millisecond))
	fmt.Println(ExecuteWithTimeout(jobs, 700*time.Millisecond))
}

func ExecuteWithTimeout(jobs []Job, timeout time.Duration) (vals []int) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Execute(ctx, jobs)
}

func Execute(ctx context.Context, jobs []Job) (vals []int) {
	var wg sync.WaitGroup
	ch := make(chan int, len(jobs))

	for i := range jobs {
		job := jobs[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- job.GetVal()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-ch:
			if !ok {
				return
			}
			vals = append(vals, v)
		}
	}
}

type Job interface {
	GetVal() int
}

type SleepJob struct {
	Duration time.Duration
	Val      int
}

func (j SleepJob) GetVal() int {
	time.Sleep(j.Duration)
	return j.Val
}
