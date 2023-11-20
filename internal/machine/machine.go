package machine

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.five07.dev/go-fsm/internal/state"
	"go.five07.dev/go-fsm/internal/types"
)

type Machine struct {
	config        types.Config
	subscriptions map[string]map[uuid.UUID]func(e types.Event)
}

// NewMachine creates a new machine instance from a config
func NewMachine(cfg types.Config) *Machine {
	cfg.Init()

	m := &Machine{}
	m.config = cfg
	m.subscriptions = map[string]map[uuid.UUID]func(e types.Event){}

	return m
}

// Config return the private config
func (m *Machine) Config() types.Config {
	return m.config
}

func (m *Machine) SetGuard(key string, guard func(ectx types.Context) bool) {
	m.config.Guards[key] = guard
}

// Subscribe
func (m *Machine) Subscribe(event string, callback func(e types.Event)) (uuid.UUID, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, err
	}

	if len(m.subscriptions[event]) == 0 {
		m.subscriptions[event] = map[uuid.UUID]func(e types.Event){}
	}

	m.subscriptions[event][u] = callback

	return u, nil
}

func (m *Machine) Unsubscribe(event string, subscriptionID uuid.UUID) {
	if len(m.subscriptions[event]) > 0 {
		delete(m.subscriptions[event], subscriptionID)
	}
}

func (m *Machine) Transistion(event string, currentState *state.State) (state.State, error) {
	var stateConfig types.StateConfig

	// this is the starting context, it can be mutated here?
	ctx := currentState.Context()
	eventTokens := make([]string, 0)

	// this only works for top-level states
	// perhaps nested states, we express the value in dot notation
	// state.Value() == parent.child.grandchild
	for key, config := range m.config.States {
		if key == currentState.Value() {
			stateConfig = config

			// append a token as we traverse into children
			eventTokens = append(eventTokens, currentState.Value())
		}
	}

	// TODO - check for parallel states, spawn parallel transitions and return

	// dot notation
	if strings.Contains(event, ".") {
		tokens := strings.Split(event, ".")
		eventTokens = tokens

		for i, token := range tokens {
			if i == len(tokens)-1 {
				event = token
			} else {
				// TODO - eventually this needs to drill into states within the element
				// right now, there are no children states in elements
				stateConfig = m.config.States[token]
			}
		}
	} else {
		eventTokens = append(eventTokens, event)
	}

	if len(stateConfig.Events) == 0 {
		return state.State{}, fmt.Errorf("the event: %v does not exist", event)
	}

	stateEventConfig := stateConfig.Events[event]

	if stateEventConfig == (types.StateEventConfig{}) {
		return state.State{}, fmt.Errorf("the event: %v does not exist", event)
	}

	isValid := true
	if stateEventConfig.Guard != "" {
		guard := m.Config().Guards[stateEventConfig.Guard]

		if guard != nil {
			isValid = guard(ctx)
		}
	}

	var newState *state.State
	event = strings.Join(eventTokens, ".")

	if isValid {
		// create a dot notated event name
		newState = state.NewState(stateEventConfig.Target, ctx)
	} else {
		newState = currentState
	}

	// check for subsciptions
	if len(m.subscriptions[event]) > 0 {
		// TODO: enqueue these, respect timer

		var wg sync.WaitGroup
		for _, callback := range m.subscriptions[event] {
			wg.Add(1)

			go func(ee types.StateEventConfig, cb func(e types.Event)) {
				defer wg.Done()

				if ee.Delay.Duration.String() != "0s" {
					time.Sleep(ee.Delay.Duration)
				}

				e := types.Event{
					Name:          event,
					State:         *newState,
					PreviousState: *currentState,
				}
				cb(e)
			}(stateEventConfig, callback)
		}

		/**
		This blocks until delayed events complete
		I think we should enqueue event, and the queue should
		respect the delay but also not block the main thread.
		*/
		wg.Wait()
	}

	// do we need a error for guards?

	return *newState, nil
}
