package config

type Cache struct {
	Expiration      int `yaml:"expiration"`
	CleanupInterval int `yaml:"cleanup_interval"`
}
