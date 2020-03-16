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
	e.set("config", argusConfigPath)
	e.set("resource", argusResourcePath)
	e.set("log", argusLogPath)
	e.set("mode", argusMode)
	return e
}

func (e Env) Get(key string) string {
	return e[key]
}

func (e Env) set(key string, value string) {
	e[key] = value
}
