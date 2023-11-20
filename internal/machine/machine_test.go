package machine

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.five07.dev/go-fsm/internal/config"
	"go.five07.dev/go-fsm/internal/state"
	"go.five07.dev/go-fsm/internal/types"
)

func TestMachine(t *testing.T) {
	bytes, err := os.ReadFile("../../examples/switch.json")
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.LoadConfig(bytes)
	assert.NoError(t, err)

	t.Run("TestNewMachine", func(t *testing.T) {
		machine := NewMachine(cfg)
		retCfg := machine.Config()

		assert.NotNil(t, machine)
		assert.Equal(t, cfg, retCfg)
	})

	t.Run("TestSubscribe", func(t *testing.T) {
		machine := NewMachine(cfg)
		subscriptionID, err := machine.Subscribe("TEST", func(e types.Event) {})

		assert.NoError(t, err)
		assert.IsType(t, uuid.Nil, subscriptionID)
		assert.Len(t, machine.subscriptions, 1)
		assert.Len(t, machine.subscriptions["TEST"], 1)
	})

	t.Run("TestUnsubscribe", func(t *testing.T) {
		machine := NewMachine(cfg)
		subscriptionID, err := machine.Subscribe("TEST", func(e types.Event) {})
		machine.Unsubscribe("OTHER", subscriptionID)

		assert.NoError(t, err)
		assert.Len(t, machine.subscriptions, 1)
		assert.Len(t, machine.subscriptions["OTHER"], 0)
		assert.Len(t, machine.subscriptions["TEST"], 1)

		machine.Unsubscribe("TEST", subscriptionID)
		assert.Len(t, machine.subscriptions, 1)
		assert.Len(t, machine.subscriptions["TEST"], 0)
	})

	t.Run("TestTransition", func(t *testing.T) {
		machine := NewMachine(cfg)
		_, err := machine.Subscribe("OFF.SWITCH", func(e types.Event) {
			assert.Equal(t, "OFF.SWITCH", e.Name)
			assert.Equal(t, "ON", e.State.Value())
		})
		machine.Transistion("SWITCH", state.NewState("OFF", map[string]interface{}{}))
		machine.Transistion("OFF.SWITCH", state.NewState("OFF", map[string]interface{}{}))

		assert.NoError(t, err)
	})
}
