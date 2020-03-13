package argus

import "os"

type Env map[string]string

var EnvVars Env

func NewEnv() Env {
	e := make(Env)
	argusResourcePath := os.Getenv("ARGUS_RESOURCE_PATH")
	argusLogPath := os.Getenv("ARGUS_LOG_PATH")
	e.set("resource", argusResourcePath)
	e.set("log", argusLogPath)
	return e
}

func (e *Env) Get(key string) string {
	return (*e)[key]
}

func (e *Env) set(key string, value string) {
	(*e)[key] = value
}
