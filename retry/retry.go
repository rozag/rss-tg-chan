package retry

import (
	"time"
)

// Policy specifies retry policy (max number of retries, sleep duration and backoff multiplier)
type Policy struct {
	Retries int
	Sleep   time.Duration
	Backoff uint
}

func do(retries int, sleep time.Duration, backoff uint, job func() error) error {
	if retries == 0 {
		return job()
	}
	err := job()
	if err == nil {
		return nil
	}
	time.Sleep(sleep)
	return do(retries-1, time.Duration(backoff)*sleep, backoff, job)
}

// Do tries to execute a job (func() error) with a number of retries, a sleep and a backoff multiplier
func Do(policy *Policy, job func() error) error {
	var sleep time.Duration
	if policy.Sleep == 0 {
		sleep = time.Second
	} else {
		sleep = policy.Sleep
	}

	var backoff uint
	if policy.Backoff == 0 {
		backoff = 2
	} else {
		backoff = policy.Backoff
	}

	return do(policy.Retries, sleep, backoff, job)
}
