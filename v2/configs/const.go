package configs

const (
	// APP application name
	APP = "binlogo"
)

// Env use for ENV
type Env string

const (
	// ENV_PRO production environment
	ENV_PRO = "production"
	// ENV_DEV dev environment
	ENV_DEV = "dev"
	// ENV_TEST test environment
	ENV_TEST = "test"
)

const (
	// CONSOLE_LISTEN default value of console listen ip
	CONSOLE_LISTEN = "0.0.0.0"
	// CONSOLE_PORT default value of console listen port
	CONSOLE_PORT = "9999"
	// CLUSTER_NAME default value of cluster name
	CLUSTER_NAME = "cluster"
)
