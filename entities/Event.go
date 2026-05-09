package entities

import (
	"chain_simulation/entities/types"
	"time"
)

type Event struct {
	StartTime time.Duration
	Action    types.ActionType
	Handler   func() error
}
