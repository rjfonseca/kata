package kata

type Config struct {
	Runner struct {
		Name string `toml:"name"`
	} `toml:"runner"`
}
