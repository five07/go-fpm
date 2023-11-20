package config

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	bytes, err := os.ReadFile("../../examples/switch.json")
	if err != nil {
		log.Fatal(err)
	}

	t.Run("TestLoadConfig", func(t *testing.T) {
		cfg, err := LoadConfig(bytes)

		assert.NoError(t, err)
		assert.Equal(t, "switch", cfg.ID)
		assert.Equal(t, "OFF", cfg.Initial)
	})
}
