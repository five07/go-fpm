package actor

import (
	"github.com/google/uuid"
	"go.five07.dev/go-fsm/internal/machine"
	"go.five07.dev/go-fsm/internal/state"
	"go.five07.dev/go-fsm/internal/types"
)

type Actor struct {
	machine *machine.Machine
	context map[string]interface{}
	state   *state.State
}

// NewActor
func NewActor(machine *machine.Machine) *Actor {
	cfg := machine.Config()

	return &Actor{
		machine: machine,
		context: cfg.Context,
		state:   &state.State{},
	}
}

// Start with the initial state
func (a *Actor) Start() *Actor {
	a.state = state.NewState(a.machine.Config().Initial, a.context)

	return a
}

// StartWithState
func (a *Actor) StartWithState(state state.State) *Actor {
	a.state = &state
	a.context = state.Context()

	return a
}

// StartWithStateJson
func (a *Actor) StartWithStateJson(json string) *Actor {
	s, _ := state.NewStateFromJson(json)
	a.state = s
	a.context = s.Context()

	return a
}

// Subscribe to an event which provides a subscriptionID
func (a *Actor) Subscribe(event string, callback func(e types.Event)) (uuid.UUID, error) {
	u, err := a.machine.Subscribe(event, callback)

	return u, err
}

// Unsubscribe from an event with a subscriptionID
func (a *Actor) Unsubscribe(event string, subscriptionID uuid.UUID) {
	a.machine.Unsubscribe(event, subscriptionID)
}

// Dispatch fires an event from a local value or global dot notated expression
func (a *Actor) Dispatch(event string) (state.State, error) {
	state, err := a.machine.Transistion(event, a.state)

	if err == nil {
		a.state = &state
	}

	return state, err
}

// Context returns the current context of the actor
func (a *Actor) Context() map[string]interface{} {
	return a.context
}

// SetContext sets the current context of the actor
func (a *Actor) SetContext(ctx map[string]interface{}) *Actor {
	a.context = ctx
	a.state.SetContext(ctx)
	return a
}

// MergeContext sets the current context of the actor
func (a *Actor) MergeContext(ctx map[string]interface{}) *Actor {
	for k, v := range ctx {
		a.context[k] = v
	}
	a.state.SetContext(ctx)

	return a
}

func (a *Actor) SetContextKey(key string, value interface{}) *Actor {
	a.context[key] = value
	a.state.SetContext(a.context)

	return a
}

func (a *Actor) GetContextKey(key string) interface{} {
	return a.context[key]
}

func (a *Actor) GetContextKeyWithDefault(key string, defaultValue interface{}) interface{} {
	val := a.context[key]

	if val != nil {
		return val
	}

	return defaultValue
}

// State returns the current state of the actor
func (a *Actor) State() state.State {
	return *a.state
}
