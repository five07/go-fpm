package types

import "go.five07.dev/go-fsm/internal/state"

type Event struct {
	Name          string
	State         state.State
	PreviousState state.State
}
