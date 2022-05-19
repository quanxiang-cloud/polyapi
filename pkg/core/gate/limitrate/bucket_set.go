package limitrate

/*
import (
	"sync"
	"time"
)

const (
	buketsetDur = defaultDur * 8
)

// NewBucketSet create and init a BucketSet object.
func NewBucketSet(ratePerSecond int) *BucketSet {
	return nil
}

// BucketSet is a set of Bucket
type BucketSet struct {
	running []*Bucket
	stoped  []*Bucket

	lock sync.RWMutex
	tick *time.Ticker
}

// NewBucket create and init a Bucket object.
func (s *BucketSet) NewBucket(ratePerSecond int) *Bucket {
	return nil
}

// RemoveBucket remove a Bucket object from BucketSet.
func (s *BucketSet) RemoveBucket(b *Bucket) {
}

func (s *BucketSet) onTick() {
	for _, v := range s.running {
		v.onTick()
	}
}

// Start start a token bucketset as running
func (s *BucketSet) Start() error {
	return nil
}

// Stop stop a token bucketset.
func (s *BucketSet) Stop() {

}
*/
