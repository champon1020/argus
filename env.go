package argus

import "os"

// EnvMap is the environment variable mapping type.
type EnvMap map[string]string

// Get returns the environment variable.
func (e EnvMap) Get(key string) string {
	return e[key]
}

func (e *EnvMap) set(key string, value string) {
	(*e)[key] = value
}

// NewEnv inithialize the struct of Env.
func NewEnv() *EnvMap {
	e := make(EnvMap)
	e.set("resource", os.Getenv("ARGUS_RESOURCE_PATH"))
	e.set("mode", os.Getenv("ARGUS_MODE"))
	e.set("secret", os.Getenv("ARGUS_SECRET_PATH"))
	e.set("admin", os.Getenv("ARGUS_ADMIN_PATH"))
	e.set("pubkey", os.Getenv("ARGUS_PUBKEY_PATH"))
	e.set("dbUser", os.Getenv("ARGUS_DB_USER"))
	e.set("dbPass", os.Getenv("ARGUS_DB_PASS"))
	e.set("dbPort", os.Getenv("ARGUS_DB_PORT"))
	e.set("dbHost", os.Getenv("ARGUS_DB_HOST"))
	e.set("dbName", os.Getenv("ARGUS_DB_NAME"))
	return &e
}
