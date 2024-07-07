package local_limit

import (
	"errors"
	"sync"
	"time"
)

// LocalLimit is a instance local request limiter.
type drop struct{}

type LocalLimit struct {
	// TargetRate is the maximum desired request rate per second
	TargetRate float64

	Drops chan drop

	DropTimeout time.Duration
	MaxDrops    int64

	SleepInterval    time.Duration
	sleepIntervalSet bool
	dropStarted      bool
	stop             bool
	overflow         float64
	dropStartOnce    sync.Once
	lastIter         time.Time
}

func (l *LocalLimit) run() {
	l.dropStarted = true
	l.lastIter = time.Now()

	var iterTime time.Duration
	var drops float64
	var fillDrops int64
	var i int64

	for !l.stop {
		time.Sleep(l.SleepInterval)
		iterTime = time.Since(l.lastIter)
		drops = iterTime.Seconds()*l.TargetRate + l.overflow

		fillDrops = min(l.MaxDrops-int64(len(l.Drops)), int64(drops))

		for i = 0; i < fillDrops; i++ {
			l.Drops <- drop{}
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

func (l *LocalLimit) Stop() {
	l.stop = true
}

func (l *LocalLimit) Limit() bool {
	if !l.dropStarted {
		l.dropStartOnce.Do(func() { go l.run() })
	}

	dropped := false

	select {
	case <-l.Drops:
		dropped = true
	case <-time.After(l.DropTimeout):
	}

	return dropped
}

func WithSleepInterval(i time.Duration) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		if i <= 0 {
			return errors.New("sleep time must be greater than 0")
		}

		l.SleepInterval = i
		l.sleepIntervalSet = true

		return nil
	}
}

func WithDDropTimeout(d time.Duration) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		l.DropTimeout = max(0*time.Second, d)

		return nil
	}
}

func WithTargetRate(r float64) func(l *LocalLimit) error {
	return func(l *LocalLimit) error {
		if r <= 0 {
			return errors.New("rate must be greater than 0")
		}

		l.TargetRate = r

		return nil
	}
}

func New(options ...func(*LocalLimit) error) (*LocalLimit, error) {
	l := LocalLimit{
		SleepInterval: 100 * time.Millisecond,
		DropTimeout:   150 * time.Millisecond,
		Drops:         make(chan drop),
		MaxDrops:      1_000_000,
		dropStarted:   false,
	}

	for _, opt := range options {
		if err := opt(&l); err != nil {
			return nil, err
		}
	}

	return &l, nil
}
