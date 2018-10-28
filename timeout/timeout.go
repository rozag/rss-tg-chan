package timeout

import (
	"errors"
	"time"
)

// Do runs the function with the timeout
func Do(timeout time.Duration, job func() error) error {
	errchan := make(chan error, 1)
	go func() { errchan <- job() }()
	select {
	case err := <-errchan:
		return err
	case <-time.After(timeout):
		return errors.New("Timed out")
	}
}
