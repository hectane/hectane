package exec

// Configuration for application execution.
type Config struct {
	Filename string `json:"-"`
	Service  bool   `json:"service"`
}
