package queue

// See https://github.com/Freeaqingme/dkim
type DKIMConfig struct {
	PrivateKey       string `json:"private-key"`
	Selector         string `json:"selector"`
	Canonicalization string `json:"canonicalization"`
}

// Application configuration.
type Config struct {
	Hostname               string `json:"hostname"`
	Directory              string `json:"directory"`
	DisableSSLVerification bool   `json:"disable-ssl-verification"`

	// Map domain names to DKIM config for that domain
	DKIMConfigs map[string]DKIMConfig `json:"dkim-configs"`
	// ProcessFunc allow you to define custom process function for message.
	ProcessFunc ProcessFunc `json:"-"`
}
