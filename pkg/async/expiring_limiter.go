package async

import (
	"time"
)

type ExpiringLimiter struct {
	*Limiter

	stopc  chan struct{}
	ticker *time.Ticker
}

func NewExpiringLimiter(limit int64, interval time.Duration) *ExpiringLimiter {
	return &ExpiringLimiter{
		Limiter: NewLimiter(limit),
		ticker:  time.NewTicker(interval),
	}
}

func (limiter *ExpiringLimiter) Start() {
	for {
		select {
		case <-limiter.ticker.C:
			limiter.Reset()
		}
	}
}

func (limiter *ExpiringLimiter) Stop() {
	limiter.stopc <- struct{}{}
}
