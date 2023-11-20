package actor_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.five07.dev/go-fsm/actor"
	"go.five07.dev/go-fsm/internal/types"
	"go.five07.dev/go-fsm/machine"
)

const stateJSON = `{
	"value": "ON",
	"context": {
		"foo": "bar"
	}
}`

func TestActor(t *testing.T) {
	bytes, err := os.ReadFile("../examples/switch.json")
	if err != nil {
		log.Fatal(err)
	}

	machine, err := machine.NewMachine(bytes)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("TestNewActor", func(t *testing.T) {
		a := actor.NewActor(machine)
		assert.NotNil(t, a)
	})

	t.Run("TestActorInitalState", func(t *testing.T) {
		a := actor.NewActor(machine)
		state := a.Start().State()
		assert.Equal(t, "OFF", state.Value())
	})

	t.Run("TestActorExistingState", func(t *testing.T) {
		a := actor.NewActor(machine)
		state := a.Start().State()
		b := actor.NewActor(machine).StartWithState(state)
		state = b.State()
		assert.Equal(t, "OFF", state.Value())
	})

	t.Run("TestActorContext", func(t *testing.T) {
		a := actor.NewActor(machine)
		ctx := a.StartWithStateJson(stateJSON).Context()
		assert.Equal(t, "bar", ctx["foo"])
	})

	t.Run("TestActorUpdateContext", func(t *testing.T) {
		a := actor.NewActor(machine)
		ctx := a.StartWithStateJson(stateJSON).Context()
		assert.Equal(t, "bar", ctx["foo"])

		ctx["foo"] = "baz"
		ctx["test"] = 123
		a.SetContext(ctx)
		assert.Equal(t, "baz", ctx["foo"])
		assert.Equal(t, 123, ctx["test"])
	})

	t.Run("TestActorMergeContext", func(t *testing.T) {
		a := actor.NewActor(machine)
		ctx := a.StartWithStateJson(stateJSON).Context()
		assert.Equal(t, "bar", ctx["foo"])

		ctx["test"] = 123
		a.MergeContext(ctx)
		assert.Equal(t, "bar", ctx["foo"])
		assert.Equal(t, 123, ctx["test"])
	})

	t.Run("TestActorDispatch", func(t *testing.T) {
		err := actor.NewActor(machine).Start().Dispatch("test")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test does not exist")

		err = actor.NewActor(machine).Start().Dispatch("SWITCH")
		assert.NoError(t, err)
	})

	t.Run("TestActorSubscribe", func(t *testing.T) {
		dispatcher := actor.NewActor(machine).Start()
		subscriber := actor.NewActor(machine).Start()

		_, err := subscriber.Subscribe("ON.SWITCH", func(e types.Event) {
			assert.Equal(t, "ON.SWITCH", e.Name)
		})
		assert.NoError(t, err)

		dispatcher.Dispatch("ON.SWITCH")
	})

	t.Run("TestActorUnsubscribe", func(t *testing.T) {
		dispatcher := actor.NewActor(machine).Start()
		subscriber := actor.NewActor(machine).Start()

		subscriptionID, err := subscriber.Subscribe("ON.SWITCH", func(e types.Event) {
			assert.Empty(t, e.Name, "This should not be called")
		})
		assert.NoError(t, err)
		subscriber.Unsubscribe("ON.SWITCH", subscriptionID)

		dispatcher.Dispatch("ON.SWITCH")
	})
}
