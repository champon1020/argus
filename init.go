package argus

var (
	// Env contains environment variables.
	Env *EnvMap

	// Logger handles application log processes.
	Logger *LogHandler
)

// Init initializes global instances.
func Init() {
	// Initialize Logger.
	Logger = NewLogger()

	// Initialize Env.
	Env = NewEnv()
}
