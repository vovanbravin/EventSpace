package auth

type Config struct {
	Database DatabaseConfig `toml:"database"`
	JWT      JWTConfig      `toml:"jwt"`
}

type DatabaseConfig struct {
	Addr string `toml:"addr"`
}

type JWTConfig struct {
	Secret string `toml:"secret"`
}
