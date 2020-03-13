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
	isTravis := os.Getenv("IS_TRAVIS")
	e.set("config", argusConfigPath)
	e.set("resource", argusResourcePath)
	e.set("log", argusLogPath)
	e.set("travis", isTravis)
	return e
}

func (e Env) Get(key string) string {
	return e[key]
}

func (e Env) set(key string, value string) {
	e[key] = value
}
