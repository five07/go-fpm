package config

import (
	"encoding/json"

	"go.five07.dev/go-fsm/internal/types"
)

func LoadConfig(machineConfig []byte) (types.Config, error) {
	config := &types.Config{}
	err := json.Unmarshal(machineConfig, config)

	return *config, err
}
