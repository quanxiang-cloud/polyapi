package limitrate

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	defaultDur = time.Millisecond
	maxDur     = defaultDur * 32
)

// TrimDuration limit duration into range
func TrimDuration(dur time.Duration) time.Duration {
	d := dur
	if d < defaultDur {
		d = defaultDur
	}
	d = d / defaultDur * defaultDur
	return d
}

func newTokenTicker(dur time.Duration, onTick func()) *tokenTicker {
	dur = TrimDuration(dur)
	n := &tokenTicker{
		ticker: nil,
		dur:    int64(dur),

		counter: 0,
		running: 0,
		onTick:  onTick,
	}
	return n
}

type tokenTicker struct {
	counter uint64

	ticker  *time.Ticker
	dur     int64
	running int32 // running status
	onTick  func()
}

func (t *tokenTicker) Start() error {
	for { // spin loop to set running => 1 when running==0
		running := atomic.LoadInt32(&t.running)
		if running > 0 || t.ticker != nil {
			return errors.New("Start on a running Bucket")
		}

		if atomic.CompareAndSwapInt32(&t.running, running, 1) {
			t.counter = 0
			break
		}
	}

	go func() {
		t.ticker = time.NewTicker(t.GetDuration())
		defer func() {
			t.ticker.Stop()
			t.ticker = nil
		}()

		for {
			select {
			case <-t.ticker.C:
				if atomic.LoadInt32(&t.running) == 0 { // stopped
					break
				}
				t.counter++
				if t.onTick != nil {
					t.onTick()
				}
			}
		}
	}()
	return nil
}

func (t *tokenTicker) Reset(dur time.Duration) {
	atomic.StoreInt64(&t.dur, int64(dur))
	if t.IsRunning() && t.ticker != nil {
		t.ticker.Reset(dur)
	}
}

func (t *tokenTicker) Stop() {
	for { // spin loop to set running => 0 when running>0
		old := atomic.LoadInt32(&t.running)
		if old <= 0 {
			break
		}
		update := int32(0)
		if atomic.CompareAndSwapInt32(&t.running, old, update) {
			if t.ticker != nil {
				t.ticker.Reset(1)
			}
			break
		}
	}
}

func (t *tokenTicker) IsRunning() bool {
	return atomic.LoadInt32(&t.running) > 0
}

func (t *tokenTicker) GetDuration() time.Duration {
	return time.Duration(atomic.LoadInt64(&t.dur))
}
