package retry

import (
	"errors"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	cases := []struct {
		policy   *Policy
		exitOn   int
		expected int
	}{
		{policy: &Policy{Retries: 0, Sleep: time.Microsecond, Backoff: 2}, exitOn: 1, expected: 1},
		{policy: &Policy{Retries: 3, Sleep: time.Microsecond, Backoff: 2}, exitOn: 4, expected: 4},
		{policy: &Policy{Retries: 0, Sleep: time.Microsecond, Backoff: 2}, exitOn: 2, expected: 1},
		{policy: &Policy{Retries: 3, Sleep: time.Microsecond, Backoff: 2}, exitOn: 2, expected: 2},
	}
	for _, c := range cases {
		c := c
		counter := 0
		Do(c.policy, func() error {
			if counter == c.exitOn {
				return nil
			}
			counter++
			return errors.New("Failed")
		})
		if counter != c.expected {
			t.Errorf("retries=%d backoff=%d exitOn=%d counter=%d expected=%d", c.policy.Retries, c.policy.Backoff, c.exitOn, counter, c.expected)
		}
	}
}
