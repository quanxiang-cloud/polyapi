package stat

import (
	"fmt"
	"sync/atomic"
	"time"
)

const averageN = 10
const unit = time.Millisecond

// NewTimeStat create a TimeStat object
func NewTimeStat(name string) *TimeStat {
	return &TimeStat{name: name}
}

// TimeStat object
type TimeStat struct {
	name    string
	min     uint64
	max     uint64
	cnt     uint64
	counter uint64
	total   uint64
}

// Add add a stat element
func (s *TimeStat) Add(dur time.Duration) {
	atomic.AddUint64(&s.counter, 1)
	if d := uint64((dur + unit/2) / unit); d > 0 {
		if min := atomic.LoadUint64(&s.min); min == 0 || d < min {
			atomic.StoreUint64(&s.min, d)
		}
		if max := atomic.LoadUint64(&s.max); d > max {
			atomic.StoreUint64(&s.max, d)
		}
		atomic.AddUint64(&s.cnt, 1)
		atomic.AddUint64(&s.total, d)
	}
}

// Clean clean the last stat info
func (s *TimeStat) Clean() {
	atomic.StoreUint64(&s.min, 0)
	atomic.StoreUint64(&s.max, 0)
	atomic.StoreUint64(&s.cnt, 0)
	atomic.StoreUint64(&s.total, 0)
}

// Report show the latest stat result
func (s *TimeStat) Report() string {
	min, max, avg, cnt, totalCnt := s.Info()
	return fmt.Sprintf("{name=%s, min=%s, max=%s, avg=%s, cnt=%d/%d}", s.name, min, max, avg, cnt, totalCnt)
}

// Info report the time stat info
func (s *TimeStat) Info() (min, max, avg time.Duration, cnt, totalCnt uint64) {
	min = time.Duration(atomic.LoadUint64(&s.min))
	max = time.Duration(atomic.LoadUint64(&s.max))
	cnt = atomic.LoadUint64(&s.cnt)
	avg = time.Duration(atomic.LoadUint64(&s.total))
	counter := atomic.LoadUint64(&s.counter)
	min *= unit
	max *= unit
	avg *= unit
	if cnt > 0 {
		avg /= time.Duration(cnt)
	}

	return min, max, avg, cnt, counter
}

// Name get name of the stat object
func (s *TimeStat) Name() string {
	return s.name
}

// Average get the average stat time
func (s *TimeStat) Average() (avg time.Duration) {
	cnt := atomic.LoadUint64(&s.cnt)
	avg = time.Duration(atomic.LoadUint64(&s.total)) * unit
	if cnt > 0 {
		avg /= time.Duration(cnt)
	}
	return
}
