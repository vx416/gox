package cache

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/cenk/backoff"
	"github.com/go-redis/redis/v8"
)

// Locker redis distribured lock
type Locker interface {
	Lock(ctx context.Context, key string, ttls ...time.Duration) (Releaser, error)
}

// Releaser represent releasable lock
type Releaser interface {
	Release(ctx context.Context) error
}

// New construct a lock
func NewLocker(client redis.Cmdable, prefix string, defaultTTL time.Duration) Locker {
	expBackOff := newExpBackOff()
	opts := &redislock.Options{
		RetryStrategy: expBackOff,
	}
	if defaultTTL <= 0 {
		defaultTTL = time.Second
	}

	return &locker{
		prefix:     prefix,
		lockClient: redislock.New(client),
		opts:       opts,
		backOff:    expBackOff,
		defaultTTL: defaultTTL,
	}
}

type locker struct {
	prefix     string
	lockClient *redislock.Client
	opts       *redislock.Options
	backOff    *expBackOff
	defaultTTL time.Duration
}

func (locker *locker) Lock(ctx context.Context, key string, ttls ...time.Duration) (Releaser, error) {
	ttl := locker.defaultTTL
	if len(ttls) > 0 {
		ttl = ttls[0]
	}
	lockKey := locker.prefix + "." + key

	defer locker.backOff.Reset()
	lock, err := locker.lockClient.Obtain(ctx, lockKey, ttl, locker.opts)
	if err != nil {
		return lock, err
	}

	return lock, nil
}

func newExpBackOff() *expBackOff {
	expBackOff := &expBackOff{backoff.NewExponentialBackOff()}
	expBackOff.backOff.InitialInterval = 10 * time.Millisecond
	expBackOff.backOff.MaxInterval = 300 * time.Millisecond
	expBackOff.backOff.MaxElapsedTime = time.Second
	expBackOff.backOff.Multiplier = 2
	return expBackOff
}

type expBackOff struct {
	backOff *backoff.ExponentialBackOff
}

func (exp *expBackOff) Reset() {
	exp.backOff.Reset()
}

func (exp *expBackOff) NextBackoff() time.Duration {
	return exp.backOff.NextBackOff()
}
