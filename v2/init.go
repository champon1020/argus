package argus

var (
	// Env contains environment variables.
	Env *EnvMap

	// Config contains application configuration.
	Config *Configuration

	// Logger handles application log processes.
	Logger *LogHandler
)

// Init initializes global instances.
func Init() {
	// Initialize Logger.
	Logger = NewLogger()

	// Initialize Env.
	Env = NewEnv()

	// Initialize Config.
	Config = NewConfig()
}
