// Strategies influenced by findings from
// https://www.awsarchitectureblog.com/2015/03/backoff.html

package time

import (
	"math/rand"
	"time"
)

// IntervalStrategy defines an interval strategy function.
type IntervalStrategy func(prev time.Duration) (next time.Duration)

// Ticker is similar to package time's Ticker but instead of choosing tick
// interval durations, tick interval strategies and maximum value are chosen.
type Ticker struct {
	C      chan time.Time
	strat  IntervalStrategy
	prev   time.Duration
	done   bool
	doneCh chan struct{}
}

// NewTicker returns a new Ticker containing a channel that will send the time
// in interval periods specified by an interval strategy. Stop the ticker to
// release associated resources.
func NewTicker(strategy IntervalStrategy) *Ticker {
	c := make(chan time.Time, 1)

	t := &Ticker{
		C:     c,
		strat: strategy,
	}

	go t.start()
	return t
}

func (t *Ticker) start() {
	t.prev = t.strat(t.prev)
	v := <-time.After(t.prev)
	if !t.done {
		t.C <- v
		go t.start()
	}
}

// Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does
// not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t *Ticker) Stop() {
	t.done = true
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only. While Tick is useful for clients that have no need to shut down
// the Ticker, be aware that without a way to shut it down the underlying Ticker
// cannot be recovered by the garbage collector; it "leaks".
func Tick(strategy IntervalStrategy) <-chan time.Time {
	return NewTicker(strategy).C
}

// Linear defines the strategy of having constant tick intervals.
func Linear(prev time.Duration) time.Duration {
	return prev
}

// Exponential defines the strategy of having linearly increasing tick
// intervals.
func Exponential(prev time.Duration) time.Duration {
	return 2 * prev
}

// LinearJitter defines the strategy of having random tick intervals.
func LinearJitter(prev time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(prev)))
}

// ExponentialJitter defines the strategy of having linearly increasing random
// tick intervals.
func ExponentialJitter(prev time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(2 * prev)))
}
