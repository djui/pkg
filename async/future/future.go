// Example use case:
//
// future := New(func() (interface{}, error) {
//     resp, err := http.Get("http://example.com")
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()
//     return ioutil.ReadAll(resp.Body)
// })
//
// // do something else
//
// b, err := future()
// body, _ := b.([]byte)

package future

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrCanceled indicates a task got canceled.
	ErrCanceled = errors.New("canceled")
	// ErrTimeout indicates a task timed out.
	ErrTimeout = errors.New("timeout")
)

// Task defines a computation as function with result value interface{} and
// error.
type Task func(...interface{}) (interface{}, error)

// Future holds a task and it's execution state.
type Future struct {
	task Task
	c    chan struct{}

	mu       sync.Mutex
	canceled bool
	done     bool
	val      interface{}
	err      error
}

// New creates a new future. Calling the return function will block until the
// future's result was computed. If an error occurs, the error will be returned
// and the result will be nil.
func New(t Task) *Future {
	future := &Future{c: make(chan struct{}, 1)}

	go func() {
		defer close(future.c)

		val, err := t()
		future.mu.Lock()
		defer future.mu.Unlock()
		future.val = val
		future.err = err
		future.done = true
		// Special case where upstream future was canceled
		if err == ErrCanceled {
			future.canceled = true
		}
	}()

	return future
}

// NewCompleted returns a new future that is already completed with the given
// value.
func NewCompleted(val interface{}, err error) *Future {
	return &Future{
		done: true,
		val:  val,
		err:  err,
	}
}

// IsDone indicates if the task is done.
func (f *Future) IsDone() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.done
}

// IsCanceled indicates if the task was canceled.
func (f *Future) IsCanceled() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.canceled
}

// Cancel cancels the task.
func (f *Future) Cancel() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.done {
		return false
	}

	f.done = true
	f.canceled = true
	return true
}

// Complete sets the value returned by get() and related methods to the given
// value, if not already completed.
func (f *Future) Complete(val interface{}) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.done && !f.canceled {
		f.done = true
		f.val = val
		f.err = nil
		return true
	}
	return false
}

// CompleteWithError sets the error returned by get() and related methods, if
// not already completed.
func (f *Future) CompleteWithError(err error) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.done && !f.canceled {
		f.done = true
		f.val = nil
		f.err = err
		return true
	}
	return false
}

// Get waits if necessary for the computation to complete, and then retrieves
// its result.
func (f *Future) Get() (interface{}, error) {
	f.mu.Lock()
	if f.canceled {
		defer f.mu.Unlock()
		return nil, ErrCanceled
	}
	if f.done {
		defer f.mu.Unlock()
		return f.val, f.err
	}
	f.mu.Unlock()

	<-f.c

	f.mu.Lock()
	defer f.mu.Unlock()
	return f.val, f.err
}

// GetNow returns the result value (or throws any encountered exception) if
// completed, else returns the given valueIfAbsent.
func (f *Future) GetNow(valueIfAbsent interface{}) (interface{}, error) {
	f.mu.Lock()
	if f.canceled {
		defer f.mu.Unlock()
		return nil, ErrCanceled
	}
	if f.done {
		defer f.mu.Unlock()
		return f.val, f.err
	}
	f.mu.Unlock()

	select {
	case <-f.c:
		f.mu.Lock()
		defer f.mu.Unlock()
		return f.val, f.err
	default:
		return valueIfAbsent, nil
	}
}

// GetWithTimeout waits if necessary for at most the given time for the
// computation to complete, and then retrieves its value, if available.
func (f *Future) GetWithTimeout(d time.Duration) (interface{}, error) {
	f.mu.Lock()
	if f.canceled {
		defer f.mu.Unlock()
		return nil, ErrCanceled
	}
	if f.done {
		defer f.mu.Unlock()
		return f.val, f.err
	}
	f.mu.Unlock()

	select {
	case <-f.c:
		f.mu.Lock()
		defer f.mu.Unlock()
		return f.val, f.err
	case <-time.After(d):
		return nil, ErrTimeout
	}
}

// ThenCompose chains a future f with given task t. The result of f is passed as
// parameter to g.
func (f *Future) ThenCompose(t Task) *Future {
	tf := func(...interface{}) (interface{}, error) {
		val, err := f.Get()
		if err != nil {
			return nil, err
		}

		return t(val)
	}

	return New(tf)
}

// ThenCombine passes the result of f and t to c if both complete without error.
func (f *Future) ThenCombine(t Task, c Task) *Future {
	tf := func(...interface{}) (interface{}, error) {
		var wg sync.WaitGroup
		var valF, valT interface{}
		var errF, errT error

		wg.Add(1)
		go func() {
			defer wg.Done()
			valF, errF = f.Get()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			valT, errT = t()
		}()

		wg.Wait()

		if errF != nil {
			return nil, errF
		}
		if errT != nil {
			return nil, errT
		}

		return c(valF, valT)
	}

	return New(tf)
}
