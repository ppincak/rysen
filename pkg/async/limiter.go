package async

import (
	"github.com/ppincak/rysen/api"
)

type Limiter struct {
	*Counter
	limit int64
}

func NewLimiter(limit int64) *Limiter {
	return &Limiter{
		Counter: NewCounter(),
		limit:   limit,
	}
}

func (limiter *Limiter) IncLim() error {
	limiter.Inc()
	if limiter.Value() >= limiter.limit {
		return api.NewError("Limit of [%d] exceeded", limiter.limit)
	}
	return nil
}
