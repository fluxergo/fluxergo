package voice

import (
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
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

type (
	LiveKitConnCreateFunc func() *LiveKitConn

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
)

func NewLivekitConn() *LiveKitConn {
	logger := slog.Default()

	conn := &LiveKitConn{
		logger: logger,
	}

	room := lksdk.NewRoom(&lksdk.RoomCallback{
		OnDisconnected: func() {
			logger.Info("disconnected from voice")
		},
		OnDisconnectedWithReason: func(reason lksdk.DisconnectionReason) {
			logger.Debug("disconnected from voice with reason", slog.Any("reason", reason))
		},
		OnParticipantConnected: func(participant *lksdk.RemoteParticipant) {
			logger.Debug("participant connected", slog.String("participant_id", participant.SID()), slog.String("identity", participant.Identity()))
		},
		OnParticipantDisconnected: func(participant *lksdk.RemoteParticipant) {
			logger.Debug("participant disconnected", slog.String("participant_id", participant.SID()), slog.String("identity", participant.Identity()))
		},
		OnActiveSpeakersChanged: func(participants []lksdk.Participant) {
			logger.Debug("active speakers changed", slog.Int("count", len(participants)))
		},
		OnRoomMetadataChanged: func(metadata string) {
			logger.Debug("room metadata changed", slog.String("metadata", metadata))
		},
		OnRecordingStatusChanged: func(isRecording bool) {
			logger.Debug("recording status changed", slog.Bool("is_recording", isRecording))
		},
		OnRoomMoved: func(roomName string, token string) {
			logger.Debug("room moved", slog.String("room_name", roomName))
		},
		OnReconnecting: func() {
			logger.Debug("reconnecting to voice")
		},
		OnReconnected: func() {
			logger.Debug("reconnected to voice")
		},
		OnLocalTrackSubscribed: func(publication *lksdk.LocalTrackPublication, lp *lksdk.LocalParticipant) {
			logger.Debug("local track subscribed", slog.String("track_sid", publication.SID()), slog.String("participant_id", lp.SID()))
		},
		ParticipantCallback: lksdk.ParticipantCallback{
			OnLocalTrackPublished: func(publication *lksdk.LocalTrackPublication, lp *lksdk.LocalParticipant) {
				logger.Debug("local track published", slog.String("track_sid", publication.SID()), slog.String("participant_id", lp.SID()))
			},
			OnLocalTrackUnpublished: func(publication *lksdk.LocalTrackPublication, lp *lksdk.LocalParticipant) {
				logger.Debug("local track unpublished", slog.String("track_sid", publication.SID()), slog.String("participant_id", lp.SID()))
			},
			OnTrackMuted: func(pub lksdk.TrackPublication, p lksdk.Participant) {
				logger.Debug("track muted", slog.String("track_sid", pub.SID()), slog.String("participant_id", p.SID()))
			},
			OnTrackUnmuted: func(pub lksdk.TrackPublication, p lksdk.Participant) {
				logger.Debug("track unmuted", slog.String("track_sid", pub.SID()), slog.String("participant_id", p.SID()))
			},
			OnMetadataChanged: func(oldMetadata string, p lksdk.Participant) {
				logger.Debug("metadata changed", slog.String("participant_id", p.SID()), slog.String("old_metadata", oldMetadata))
			},
			OnAttributesChanged: func(changed map[string]string, p lksdk.Participant) {
				logger.Debug("attributes changed", slog.String("participant_id", p.SID()), slog.Any("changed", changed))
			},
			OnIsSpeakingChanged: func(p lksdk.Participant) {
				logger.Debug("is speaking changed", slog.String("participant_id", p.SID()), slog.Bool("is_speaking", p.IsSpeaking()))
			},
			OnConnectionQualityChanged: func(update *livekit.ConnectionQualityInfo, p lksdk.Participant) {
				logger.Debug("connection quality changed", slog.String("participant_id", p.SID()), slog.Any("quality", update.Quality))
			},
			OnTrackSubscribed: func(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
				logger.Debug("track subscribed", slog.String("stream_id", track.StreamID()), slog.String("participant_id", rp.SID()))
			},
			OnTrackUnsubscribed: func(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
				logger.Debug("track unsubscribed", slog.String("stream_id", track.StreamID()), slog.String("participant_id", rp.SID()))
			},
			OnTrackSubscriptionFailed: func(sid string, rp *lksdk.RemoteParticipant) {
				logger.Debug("track subscription failed", slog.String("track_sid", sid), slog.String("participant_id", rp.SID()))
			},
			OnTrackPublished: func(publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
				logger.Debug("track published", slog.String("track_sid", publication.SID()), slog.String("participant_id", rp.SID()))
			},
			OnTrackUnpublished: func(publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
				logger.Debug("track unpublished", slog.String("track_sid", publication.SID()), slog.String("participant_id", rp.SID()))
			},
			OnDataReceived: func(data []byte, params lksdk.DataReceiveParams) {
				logger.Debug("data received", slog.Int("size", len(data)), slog.String("participant_id", params.SenderIdentity))
			},
			OnDataPacket: func(data lksdk.DataPacket, params lksdk.DataReceiveParams) {
				logger.Debug("data packet received", slog.String("participant_id", params.SenderIdentity))
			},
			OnTranscriptionReceived: func(transcriptionSegments []*lksdk.TranscriptionSegment, p lksdk.Participant, publication lksdk.TrackPublication) {
				logger.Debug("transcription received", slog.String("participant_id", p.SID()), slog.Int("segments", len(transcriptionSegments)))
			},
		},
	})
	conn.room = room

	return conn
}

type LiveKitConn struct {
	logger *slog.Logger
	room   *lksdk.Room
}

func (l *LiveKitConn) Open(state State) error {
	return l.room.JoinWithToken(state.Endpoint, state.Token)
}

func (l *LiveKitConn) Close() {
	l.room.Disconnect()
}

func (l *LiveKitConn) Status() Status {
	switch l.room.ConnectionState() {
	case lksdk.ConnectionStateConnected:
		return StatusConnected
	case lksdk.ConnectionStateReconnecting:
		return StatusConnecting
	default:
		return StatusDisconnected
	}
}

type AudioSource int

const (
	AudioSourceMicrophone AudioSource = iota
	AudioSourceScreenShare
)

func (s AudioSource) toLiveKit() livekit.TrackSource {
	switch s {
	case AudioSourceMicrophone:
		return livekit.TrackSource_MICROPHONE
	case AudioSourceScreenShare:
		return livekit.TrackSource_SCREEN_SHARE_AUDIO
	default:
		return livekit.TrackSource_UNKNOWN
	}
}

type VideoSource int

const (
	VideoSourceCamera VideoSource = iota
	VideoSourceScreenShare
)

func (s VideoSource) toLiveKit() livekit.TrackSource {
	switch s {
	case VideoSourceCamera:
		return livekit.TrackSource_CAMERA
	case VideoSourceScreenShare:
		return livekit.TrackSource_SCREEN_SHARE
	default:
		return livekit.TrackSource_UNKNOWN
	}
}

func (l *LiveKitConn) AudioWriter(name string, source AudioSource) (io.WriteCloser, error) {
	return l.writer(webrtc.RTPCodecCapability{
		MimeType:  webrtc.MimeTypeOpus,
		ClockRate: 48000,
		Channels:  2,
	}, lksdk.TrackPublicationOptions{
		Name:   name,
		Source: source.toLiveKit(),
	}, 0)
}

func (l *LiveKitConn) VideoWriter(name string, source VideoSource, width int, height int, fps int) (io.WriteCloser, error) {
	return l.writer(webrtc.RTPCodecCapability{
		MimeType:  webrtc.MimeTypeH264,
		ClockRate: 90000,
		Channels:  0,
	}, lksdk.TrackPublicationOptions{
		Name:        name,
		Source:      source.toLiveKit(),
		VideoWidth:  width,
		VideoHeight: height,
	}, fps)
}

func (l *LiveKitConn) writer(c webrtc.RTPCodecCapability, options lksdk.TrackPublicationOptions, fps int) (io.WriteCloser, error) {
	track, err := lksdk.NewLocalTrack(c)
	if err != nil {
		return nil, err
	}

	_, err = l.room.LocalParticipant.PublishTrack(track, &options)
	if err != nil {
		return nil, err
	}

	var sampleDuration time.Duration
	switch c.MimeType {
	case webrtc.MimeTypeOpus:
		sampleDuration = time.Millisecond * 20
	case webrtc.MimeTypeH264:
		sampleDuration = time.Second / time.Duration(fps)
	default:
		return nil, errors.New("unsupported codec")
	}

	return &trackWriter{
		track:          track,
		sampleDuration: sampleDuration,
	}, nil
}

type trackWriter struct {
	track          *lksdk.LocalTrack
	sampleDuration time.Duration
}

func (t *trackWriter) Write(p []byte) (int, error) {
	err := t.track.WriteSample(media.Sample{
		Data:     p,
		Duration: t.sampleDuration,
	}, nil)

	return len(p), err
}

func (t trackWriter) Close() error {
	return t.track.Close()
}
