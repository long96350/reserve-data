package stat

import (
	"time"
)

type FetcherRunner interface {
	GetBlockTicker() <-chan time.Time
	GetLogTicker() <-chan time.Time
	GetReserveRatesTicker() <-chan time.Time
	Start() error
	Stop() error
}