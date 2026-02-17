package voice

import (
	"log/slog"
)

func defaultConnConfig() connConfig {
	return connConfig{
		Logger:                slog.Default(),
		LiveKitConnCreateFunc: NewLivekitConn,
	}
}

type connConfig struct {
	Logger *slog.Logger

	LiveKitConnCreateFunc LiveKitConnCreateFunc
}

// ConnConfigOpt is used to functionally configure a connConfig.
type ConnConfigOpt func(config *connConfig)

func (c *connConfig) apply(opts []ConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "voice_conn"))
}

// WithConnLogger sets the Conn(s) used Logger.
func WithConnLogger(logger *slog.Logger) ConnConfigOpt {
	return func(config *connConfig) {
		config.Logger = logger
	}
}

// WithConnLiveKitConnCreateFunc sets the Conn(s) used LiveKitConnCreateFunc.
func WithConnLiveKitConnCreateFunc(liveKitConnCreateFunc LiveKitConnCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.LiveKitConnCreateFunc = liveKitConnCreateFunc
	}
}
