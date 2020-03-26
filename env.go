package argus

import (
	"os"
)

type Env map[string]string

var EnvVars Env

func NewEnv() Env {
	e := make(Env)
	argusConfigPath := os.Getenv("ARGUS_CONFIG_PATH")
	argusResourcePath := os.Getenv("ARGUS_RESOURCE_PATH")
	argusLogPath := os.Getenv("ARGUS_LOG_PATH")
	argusMode := os.Getenv("ARGUS_MODE")
	argusSecretPath := os.Getenv("ARGUS_SECRET_PATH")
	argusUserPath := os.Getenv("ARGUS_USER_PATH")
	argusKeyPath := os.Getenv("ARGUS_KEY_PATH")
	e.set("config", argusConfigPath)
	e.set("resource", argusResourcePath)
	e.set("log", argusLogPath)
	e.set("mode", argusMode)
	e.set("secret", argusSecretPath)
	e.set("user", argusUserPath)
	e.set("key", argusKeyPath)
	return e
}

func (e Env) Get(key string) string {
	return e[key]
}

func (e Env) set(key string, value string) {
	e[key] = value
}
