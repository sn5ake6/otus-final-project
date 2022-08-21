package bucket

import (
	"context"
	"sync"
	"time"

	"github.com/sn5ake6/otus-final-project/internal/config"
	"github.com/sn5ake6/otus-final-project/internal/storage"
)

type values map[string]int64

type bucket struct {
	limit  int64
	values values
}

type LeakyBucket struct {
	ctx            context.Context
	mu             *sync.Mutex
	loginBucket    bucket
	passwordBucket bucket
	ipBucket       bucket
	resetTicker    *time.Ticker
}

func NewLeakyBucket(ctx context.Context, limits config.LimitConf) *LeakyBucket {
	resetInterval, _ := time.ParseDuration(limits.ResetInterval)
	resetTicker := time.NewTicker(resetInterval)
	go func() {
		<-ctx.Done()
		resetTicker.Stop()
	}()

	lb := &LeakyBucket{
		ctx:            ctx,
		mu:             &sync.Mutex{},
		loginBucket:    newBuckets(limits.Login),
		passwordBucket: newBuckets(limits.Password),
		ipBucket:       newBuckets(limits.IP),
		resetTicker:    resetTicker,
	}

	go lb.leak()

	return lb
}

func newBuckets(limit int64) bucket {
	return bucket{
		limit:  limit,
		values: make(values),
	}
}

func (lb *LeakyBucket) Check(authorize storage.Authorize) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	return lb.checkLimit(lb.loginBucket, authorize.Login) &&
		lb.checkLimit(lb.passwordBucket, authorize.Password) &&
		lb.checkLimit(lb.ipBucket, authorize.IP)
}

func (lb *LeakyBucket) checkLimit(b bucket, value string) bool {
	if b.values[value] < b.limit {
		b.values[value]++
		return true
	}

	return false
}

func (lb *LeakyBucket) Reset(authorize storage.Authorize) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	delete(lb.loginBucket.values, authorize.Login)
	delete(lb.passwordBucket.values, authorize.Password)
	delete(lb.ipBucket.values, authorize.IP)
}

func (lb *LeakyBucket) resetAll() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.loginBucket.values = make(values)
	lb.passwordBucket.values = make(values)
	lb.ipBucket.values = make(values)
}

func (lb *LeakyBucket) leak() {
	for range lb.resetTicker.C {
		select {
		case <-lb.ctx.Done():
			return
		default:
			lb.resetAll()
		}
	}
}
