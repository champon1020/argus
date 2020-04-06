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
	argusMode := os.Getenv("ARGUS_MODE")
	argusSecretPath := os.Getenv("ARGUS_SECRET_PATH")
	argusUserPath := os.Getenv("ARGUS_USER_PATH")
	argusKeyPath := os.Getenv("ARGUS_KEY_PATH")
	argusDbUser := os.Getenv("ARGUS_DB_USER")
	argusDbPass := os.Getenv("ARGUS_DB_PASS")
	argusDbPort := os.Getenv("ARGUS_DB_PORT")
	argusDbHost := os.Getenv("ARGUS_DB_HOST")
	argusDbName := os.Getenv("ARGUS_DB_NAME")
	e.set("config", argusConfigPath)
	e.set("resource", argusResourcePath)
	e.set("mode", argusMode)
	e.set("secret", argusSecretPath)
	e.set("user", argusUserPath)
	e.set("key", argusKeyPath)
	e.set("dbUser", argusDbUser)
	e.set("dbPass", argusDbPass)
	e.set("dbPort", argusDbPort)
	e.set("dbHost", argusDbHost)
	e.set("dbName", argusDbName)
	return e
}

func (e Env) Get(key string) string {
	return e[key]
}

func (e Env) set(key string, value string) {
	e[key] = value
}
