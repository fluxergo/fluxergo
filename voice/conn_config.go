package voice

import (
	"log/slog"
)

func defaultConnConfig() connConfig {
	return connConfig{
		Logger:                  slog.Default(),
		AudioSenderCreateFunc:   NewAudioSender,
		AudioReceiverCreateFunc: NewAudioReceiver,
	}
}

type connConfig struct {
	Logger *slog.Logger

	AudioSenderCreateFunc   AudioSenderCreateFunc
	AudioReceiverCreateFunc AudioReceiverCreateFunc
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

// WithConnAudioSenderCreateFunc sets the Conn(s) used AudioSenderCreateFunc.
func WithConnAudioSenderCreateFunc(audioSenderCreateFunc AudioSenderCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.AudioSenderCreateFunc = audioSenderCreateFunc
	}
}

// WithConnAudioReceiverCreateFunc sets the Conn(s) used AudioReceiverCreateFunc.
func WithConnAudioReceiverCreateFunc(audioReceiverCreateFunc AudioReceiverCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.AudioReceiverCreateFunc = audioReceiverCreateFunc
	}
}
