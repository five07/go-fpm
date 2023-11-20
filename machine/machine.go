package machine

import (
	"go.five07.dev/go-fsm/internal/config"
	"go.five07.dev/go-fsm/internal/machine"
)

// NewMachine creates a new machine from a config
func NewMachine(machineConfig []byte) (*machine.Machine, error) {
	cfg, err := config.LoadConfig(machineConfig)
	if err != nil {
		return nil, err
	}

	machine := machine.NewMachine(cfg)

	return machine, err
}
