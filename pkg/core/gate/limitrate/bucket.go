package limitrate

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	minTokenPerSecond = 1
	maxTokenPerSecond = 100000
	maxCapacity       = 1 // tokens for N second(s)
	bucketDur         = defaultDur * 8
)

// NewBucket create and init a Bucket object.
func NewBucket(tokenPerSecond int) *Bucket {
	tokenPerSecond = trimLimit(tokenPerSecond, minTokenPerSecond, maxTokenPerSecond)
	capacity := tokenPerSecond * maxCapacity

	dur, tokenPerTick := calcTickPara(tokenPerSecond, bucketDur)

	n := &Bucket{
		tokenUse:       0,
		tokenGen:       0,
		capacity:       uint64(capacity),
		tokenPerTick:   tokenPerTick,
		tokenPerSecond: int64(tokenPerSecond),
		ticker:         nil,
		//set:            nil,
	}
	n.ticker = newTokenTicker(dur, n.onTick)

	return n
}

// Bucket is the object to limit rate by token.
type Bucket struct {
	// let tokenUse catches up with tokenGen to adjust producer and consumer
	tokenUse uint64 // token used
	tokenGen uint64 // token generated

	capacity       uint64
	tokenPerTick   uint64
	tokenPerSecond int64
	ticker         *tokenTicker
	//set            *BucketSet
}

func (b *Bucket) onTick() {
	for {
		used := atomic.LoadUint64(&b.tokenUse)
		gen := atomic.LoadUint64(&b.tokenGen)
		maxGen := used + atomic.LoadUint64(&b.capacity)
		if gen >= maxGen { //Bucket is full
			break
		}
		if gen < maxGen {
			update := gen + atomic.LoadUint64(&b.tokenPerTick)
			if update > maxGen { // do not overflow capacity
				update = maxGen
			}
			if atomic.CompareAndSwapUint64(&b.tokenGen, gen, update) {
				break
			}
		}
	}
}

// Start start a token bucket as running
func (b *Bucket) Start() error {
	// if b.set != nil {
	// 	return errors.New("Start a Bucket that is in BucketSet")
	// }
	if b.ticker == nil {
		return errors.New("Start a Bucket that without ticker")
	}
	if err := b.ticker.Start(); err != nil {
		return err
	}

	atomic.StoreUint64(&b.tokenUse, 0)
	atomic.StoreUint64(&b.tokenGen, 0)
	return nil
}

// Stop stop a running token bucket
func (b *Bucket) Stop() {
	if b.ticker != nil {
		b.ticker.Stop()
	}
}

// GetToken try to get a token from bucket.
// It return false if no token in this bucket currently.
func (b *Bucket) GetToken() bool {
	for {
		used := atomic.LoadUint64(&b.tokenUse)
		gen := atomic.LoadUint64(&b.tokenGen)
		if used >= gen { // no aviable token
			return false
		}

		update := used + 1
		if atomic.CompareAndSwapUint64(&b.tokenUse, used, update) {
			return true
		}
	}
	return false
}

// IsRunning check if a bucket is running
func (b *Bucket) IsRunning() bool {
	if b.ticker != nil {
		return b.ticker.IsRunning()
	}
	return false
}

// ShowTokenArg show token generator arg
func (b *Bucket) ShowTokenArg() string {
	if b.IsRunning() && b.ticker != nil {
		dur := b.ticker.GetDuration()
		perTick := atomic.LoadUint64(&b.tokenPerTick)
		perSec := atomic.LoadInt64(&b.tokenPerSecond)
		perSecCalc := uint64(time.Second/dur) * perTick
		used := atomic.LoadUint64(&b.tokenUse)
		gen := atomic.LoadUint64(&b.tokenGen)
		return fmt.Sprintf("TokenArg: Dur=%s TokenPerTick=%d TokenPerSecond=%d/%d state=%d/%d",
			dur, perTick, perSecCalc, perSec, used, gen)
	}
	return "TokenArg: stopped"
}

//------------------------------------------------------------------------------

func trimLimit(val, min, max int) int {
	if val < min {
		val = min
	}
	if max > 0 && val > max {
		val = max
	}
	return val
}

func calcTickPara(tokenPerSec int, minDur time.Duration) (dur time.Duration, tokenPerTick uint64) {
	dur = time.Second / time.Duration(tokenPerSec)
	if dur < minDur {
		dur = minDur
	}
	dur = TrimDuration(dur)
	return dur, _calcTokenPerTick(tokenPerSec, dur)
}

func _calcTokenPerTick(tokenPerSec int, dur time.Duration) uint64 {
	ticksPerSec := uint64((time.Second + dur - 1) / dur)
	tokenPerTick := (uint64(tokenPerSec) + ticksPerSec/2) / ticksPerSec
	if tokenPerTick <= 0 {
		tokenPerTick = 1
	}
	return tokenPerTick
}
