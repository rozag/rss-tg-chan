package retry

import (
	"errors"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	cases := []struct {
		retries  uint
		sleep    time.Duration
		backoff  uint
		exitOn   int
		expected int
	}{
		{retries: 0, sleep: time.Microsecond, backoff: 2, exitOn: 1, expected: 1},
		{retries: 3, sleep: time.Microsecond, backoff: 2, exitOn: 4, expected: 4},
		{retries: 0, sleep: time.Microsecond, backoff: 2, exitOn: 2, expected: 1},
		{retries: 3, sleep: time.Microsecond, backoff: 2, exitOn: 2, expected: 2},
	}
	for _, c := range cases {
		c := c
		counter := 0
		Do(c.retries, c.sleep, c.backoff, func() error {
			if counter == c.exitOn {
				return nil
			}
			counter++
			return errors.New("Failed")
		})
		if counter != c.expected {
			t.Errorf("retries=%d backoff=%d exitOn=%d counter=%d expected=%d", c.retries, c.backoff, c.exitOn, counter, c.expected)
		}
	}
}
