package voice

import (
	"errors"
	"io"

	"github.com/disgoorg/snowflake/v2"
)

var ErrLiveKitNotConnected = errors.New("livekit not connected")

// Status returns the current status of the gateway.
type Status int

const (
	StatusDisconnected Status = iota
	StatusConnecting
	StatusConnected
)

type State struct {
	GuildID   snowflake.ID
	UserID    snowflake.ID
	ChannelID *snowflake.ID

	SessionID    string
	Token        string
	ConnectionID string
	Endpoint     string
}

type AudioSource int

const (
	AudioSourceMicrophone AudioSource = iota
	AudioSourceScreenShare
)

type VideoSource int

const (
	VideoSourceCamera VideoSource = iota
	VideoSourceScreenShare
)

type (
	LiveKitConnCreateFunc func() LiveKitConn

	Packet struct {
		Type byte
		// Sequence is the sequence number of the packet.
		Sequence uint16
		// Timestamp is the timestamp of the packet.
		Timestamp uint32
		// SSRC is the users SSRC of the packet.
		SSRC         uint32
		HasExtension bool
		ExtensionID  int
		Extension    []byte
		CSRC         []uint32
		HeaderSize   int
		// Opus is the actual opus data of the packet.
		Opus []byte
	}

	LiveKitConn interface {
		Open(state State) error
		Close()

		Status() Status

		AudioWriter(name string, source AudioSource) (io.WriteCloser, error)
		VideoWriter(name string, source VideoSource, width int, height int, fps int) (io.WriteCloser, error)
	}
)
