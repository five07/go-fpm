package types

type Context map[string]interface{}

type Config struct {
	ID      string                 `json:"id"`
	Initial string                 `json:"initial"`
	Context Context                `json:"context"`
	States  map[string]StateConfig `json:"states"`

	Guards map[string]func(ctx Context) bool
}

func (c *Config) Init() {
	c.Guards = map[string]func(ctx Context) bool{}
}

type StateConfig struct {
	Initial string                      `json:"initial"`
	States  map[string]StateConfig      `json:"states"`
	Events  map[string]StateEventConfig `json:"events"`
}

type StateEventConfig struct {
	Target string   `json:"target"`
	Entry  string   `json:"entry"`
	Exit   string   `json:"exit"`
	Guard  string   `json:"guard"`
	Delay  Duration `json:"delay"`
}
