package state

import (
	"encoding/json"
)

type State struct {
	context map[string]interface{}
	value   string
}

type PublicState struct {
	Context map[string]interface{} `json:"context"`
	Value   string                 `json:"value"`
}

// TODO - evaluate child initial state?? person.json -> ASLEEP.WAKE_UP -> AWAKE.HOME

// NewState creates an state with value and context
func NewState(value string, context map[string]interface{}) *State {
	return &State{
		context: context,
		value:   value,
	}
}

// NewStateFromJson rehydrate JSON into a
func NewStateFromJson(data string) (*State, error) {
	s := PublicState{}

	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		return nil, err
	}

	return &State{
		context: s.Context,
		value:   s.Value,
	}, nil
}

// Value returns the state value
func (s *State) Value() string {
	return s.value
}

// Json returns the state as dehydrated JSON
func (s *State) Json() (string, error) {
	tmp := PublicState{
		Context: s.context,
		Value:   s.value,
	}

	value, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}

	return string(value[:]), nil
}

// Context returns the state's current context
func (s *State) Context() map[string]interface{} {
	return s.context
}

// SetContext sets the current context of the actor
func (s *State) SetContext(ctx map[string]interface{}) {
	s.context = ctx
}

// SetContext sets the current context of the actor
func (s *State) MergeContext(ctx map[string]interface{}) {
	for k, v := range ctx {
		s.context[k] = v
	}
}
