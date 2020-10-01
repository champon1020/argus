package argus

var (
	// Env contains environment variables.
	Env *EnvMap

	// Config contains application configuration.
	Config *Configuration

	// Errs has some errors occurred in api call.
	Errs *ErrorsHandler

	// Logger handles application log processes.
	Logger *LogHandler
)

// Init initializes global instances.
func Init() {
	// Initialize ErrorsHandler.
	Errs = new(ErrorsHandler)

	// Initialize Logger.
	Logger = NewLogger()

	// Initialize Env.
	Env = NewEnv()

	// Initialize Config.
	Config = NewConfig()
}
