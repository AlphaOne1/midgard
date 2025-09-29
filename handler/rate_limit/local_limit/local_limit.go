// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package local_limit provides a process-local rate limiter.
package local_limit

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// drop is an empty structure to be used as drops in the internal algorithm.
type drop struct{}

// LocalLimit is a instance local request limiter.
type LocalLimit struct {
	// TargetRate is the maximum desired request rate per second
	TargetRate float64

	// drops is the channel containing the drops for the requests.
	drops chan drop

	// DropTimeout is the time a request waits for a drop,
	// if there are no drops currently available in the drops channel.
	DropTimeout time.Duration

	// MaxDrops specifies the maximum number of drops in the drops channel.
	// Setting the maximum prevents services from being flooded after a longer
	// period without requests.
	MaxDrops int64

	// SleepInterval is the time between the generation of drops
	SleepInterval time.Duration
	// dropStarted signalizes if the drop generation has been started.
	dropStarted atomic.Bool
	// stop signals the internal drop generator to stop working
	stop atomic.Bool
	// overflow stores the fractional drops, especially with a low TargetRate this
	// guarantees no lost drops.
	overflow float64
	// dropStartOnce cares that the internal drop generation is just started once.
	dropStartOnce sync.Once
	// lastIter stores the last time the internal drop generator was running
	lastIter time.Time
}

// run generates the drops. It is called internally as a go routine.
func (l *LocalLimit) run() {
	l.dropStarted.Store(true)
	l.lastIter = time.Now()

	var iterTime time.Duration
	var drops float64
	var fillDrops int64

	for !l.stop.Load() {
		time.Sleep(l.SleepInterval)
		iterTime = time.Since(l.lastIter)
		drops = iterTime.Seconds()*l.TargetRate + l.overflow

		fillDrops = min(l.MaxDrops-int64(len(l.drops)), int64(drops))

		for range fillDrops {
			l.drops <- drop{}
		}

		if fillDrops == int64(drops) {
			l.overflow = drops - float64(int64(drops))
		} else {
			// if we had to cut the drops, there is also no overflow
			l.overflow = 0
		}

		l.lastIter = time.Now()
	}
}

// Stop sets the stop marker, so the drop generator can stop eventually.
func (l *LocalLimit) Stop() {
	l.stop.Store(true)
}

// Limit gives true, if the rate limit is not yet exceeded, otherwise false.
// If there are currently no drops to exhaust, it will wait the configured
// DropTimeout for a drop.
func (l *LocalLimit) Limit() bool {
	if !l.dropStarted.Load() {
		l.dropStartOnce.Do(func() { go l.run() })
	}

	dropped := false

	select {
	case <-l.drops:
		dropped = true
	case <-time.After(l.DropTimeout):
	}

	return dropped
}

// WithDropTimeout sets the timeout a process calling Limit will wait, before
// giving up to get a drop.
func WithDropTimeout(d time.Duration) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		l.DropTimeout = max(0*time.Second, d)

		return nil
	}
}

// WithMaxDropsAbsolute sets the maximum number of drops to be stored before
// dismissing them. This is to save the server from a flood of requests after
// a longer period of no requests.
func WithMaxDropsAbsolute(d int64) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		l.MaxDrops = d

		return nil
	}
}

// WithMaxDropsInterval sets the maximum number of drops to be stored before
// dismissing them. Other than WithMaxDropsAbsolute it calculates the number of
// drops to store depending on the TargetRate parameter.
// This method must be called after setting the TargetRate. Otherwise, the defaults
// are used.
func WithMaxDropsInterval(d time.Duration) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		l.MaxDrops = max(1, int64(l.TargetRate*d.Seconds()))

		return nil
	}
}

// WithSleepInterval sets the interval for the drop generator in which new drops are
// generated. New drops are generated for the last passed interval. Fractional drops
// are stored to be used in the next iteration.
func WithSleepInterval(i time.Duration) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		if i <= 0 {
			return errors.New("sleep time must be greater than 0")
		}

		l.SleepInterval = i

		return nil
	}
}

// WithTargetRate sets the requested rate of drops per second.
func WithTargetRate(r float64) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		if r <= 0 {
			return errors.New("rate must be greater than 0")
		}

		l.TargetRate = r

		return nil
	}
}

// New creates a new local rate limiter.
func New(options ...func(*LocalLimit) error) (*LocalLimit, error) {
	l := LocalLimit{
		TargetRate:    1,
		SleepInterval: 100 * time.Millisecond,
		DropTimeout:   150 * time.Millisecond,
		drops:         make(chan drop),
		MaxDrops:      1_000,
		dropStarted:   atomic.Bool{},
	}

	for _, opt := range options {
		if err := opt(&l); err != nil {
			return nil, err
		}
	}

	return &l, nil
}
