package cmd

type TelemetryConfig struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port int
}
