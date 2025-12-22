package config

type Config struct {
	UI UIConfig `toml:"ui"`
}

type UIConfig struct {
	Lang           string `toml:"lang"`
	NonInteractive bool   `toml:"non_interactive"`
}

func DefaultConfig() *Config {
	return &Config{
		UI: UIConfig{
			Lang:           "",
			NonInteractive: false,
		},
	}
}
