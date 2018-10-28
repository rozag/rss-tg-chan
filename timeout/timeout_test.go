package timeout

import (
	"fmt"
	"testing"
	"time"
)

type cs struct {
	sleep, timeout time.Duration
	expectedErr    bool
}

func (c cs) String() string {
	return fmt.Sprintf("Case{sleep=%d timeout=%d expectedErr=%v}", c.sleep, c.timeout, c.expectedErr)
}

func TestDo(t *testing.T) {
	cases := []cs{
		cs{sleep: 1 * time.Millisecond, timeout: 2 * time.Millisecond, expectedErr: false},
		cs{sleep: 2 * time.Millisecond, timeout: 1 * time.Millisecond, expectedErr: true},
		cs{sleep: 1 * time.Millisecond, timeout: 1 * time.Millisecond, expectedErr: true},
	}
	for _, c := range cases {
		c := c
		err := Do(c.timeout, func() error {
			time.Sleep(c.sleep)
			return nil
		})
		gotErr := err != nil
		switch {
		case !gotErr && c.expectedErr:
			t.Errorf("%v; Expected error but didn't get it", c)
		case gotErr && !c.expectedErr:
			t.Errorf("%v; Didn't expect error but got it", c)
		}
	}
}
