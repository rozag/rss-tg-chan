package retry

import (
	"time"
)

// Do tries to execute a job (func() error) with a number of retries, a sleep and a backoff multiplier
func Do(retries uint, sleep time.Duration, backoff uint, job func() error) error {
	if retries == 0 {
		return job()
	}
	err := job()
	if err == nil {
		return nil
	}
	time.Sleep(sleep)
	return Do(retries-1, time.Duration(backoff)*sleep, backoff, job)
}
