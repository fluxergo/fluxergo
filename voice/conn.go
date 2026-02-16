package voice

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/pion/webrtc/v4"

	"github.com/fluxergo/fluxergo/gateway"
)

type (
	// ConnCreateFunc is a type alias for a function that creates a new Conn.
	ConnCreateFunc func(guildID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn

	// Conn is a complete voice conn to fluxer. It holds the Gateway and voiceudp.UDPConn conn and combines them.
	Conn interface {
		// ChannelID returns the ID of the voice channel the voice Conn is openedChan to.
		ChannelID() *snowflake.ID

		// GuildID returns the ID of the guild the voice Conn is openedChan to.
		GuildID() snowflake.ID

		// SetOpusFrameProvider lets you inject your own OpusFrameProvider.
		SetOpusFrameProvider(handler OpusFrameProvider)

		// SetOpusFrameReceiver lets you inject your own OpusFrameReceiver.
		SetOpusFrameReceiver(handler OpusFrameReceiver)

		// Open opens the voice conn. It will connect to the voice gateway and start the Conn conn after it receives the Gateway events.
		Open(ctx context.Context, channelID snowflake.ID, selfMute bool, selfDeaf bool) error

		// Close closes the voice conn. It will close the Conn conn and disconnect from the voice gateway.
		Close(ctx context.Context)

		// HandleVoiceStateUpdate provides the gateway.EventVoiceStateUpdate to the voice conn. Which is needed to connect to the voice Gateway.
		HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate)

		// HandleVoiceServerUpdate provides the gateway.EventVoiceServerUpdate to the voice conn. Which is needed to connect to the voice Gateway.
		HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate)

		Room() *lksdk.Room
	}
)

type State struct {
	GuildID snowflake.ID
	UserID  snowflake.ID

	ChannelID    *snowflake.ID
	SessionID    string
	Token        string
	ConnectionID string
	Endpoint     string
}

// NewConn returns a new default voice conn.
func NewConn(guildID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn {
	cfg := defaultConnConfig()
	cfg.apply(opts)

	openedCtx, openedCancel := context.WithCancel(context.Background())
	closedCtx, closedCancel := context.WithCancel(context.Background())

	conn := &connImpl{
		config:               cfg,
		voiceStateUpdateFunc: voiceStateUpdateFunc,
		removeConnFunc:       removeConnFunc,
		state: State{
			GuildID: guildID,
			UserID:  userID,
		},
		openedCtx: openedCtx,
		opened:    openedCancel,
		closedCtx: closedCtx,
		closed:    closedCancel,
	}

	return conn
}

type connImpl struct {
	config               connConfig
	voiceStateUpdateFunc StateUpdateFunc
	removeConnFunc       func()

	state   State
	stateMu sync.Mutex

	room *lksdk.Room

	audioSender   AudioSender
	audioReceiver AudioReceiver

	openedCtx context.Context
	opened    context.CancelFunc
	closedCtx context.Context
	closed    context.CancelFunc
}

func (c *connImpl) ChannelID() *snowflake.ID {
	return c.state.ChannelID
}

func (c *connImpl) GuildID() snowflake.ID {
	return c.state.GuildID
}

func (c *connImpl) SetOpusFrameProvider(provider OpusFrameProvider) {
	if c.audioSender != nil {
		c.audioSender.Close()
	}
	c.audioSender = c.config.AudioSenderCreateFunc(c.config.Logger, provider, c)
	c.audioSender.Open()
}

func (c *connImpl) SetOpusFrameReceiver(handler OpusFrameReceiver) {
	if c.audioReceiver != nil {
		c.audioReceiver.Close()
	}
	c.audioReceiver = c.config.AudioReceiverCreateFunc(c.config.Logger, handler, c)
	c.audioReceiver.Open()
}

func (c *connImpl) HandleVoiceStateUpdate(update gateway.EventVoiceStateUpdate) {
	if update.GuildID != c.state.GuildID || update.UserID != c.state.UserID {
		return
	}

	if update.ChannelID == nil {
		c.state.ChannelID = nil
		c.room.Disconnect()
		c.closed()
	} else {
		c.state.ChannelID = update.ChannelID

		room := lksdk.NewRoom(&lksdk.RoomCallback{
			OnDisconnected: func() {
				c.config.Logger.Info("disconnected from voice")
			},
			OnDisconnectedWithReason: func(reason lksdk.DisconnectionReason) {
				c.config.Logger.Debug("disconnected from voice with reason", slog.Any("reason", reason))
			},
			OnParticipantConnected: func(participant *lksdk.RemoteParticipant) {
				c.config.Logger.Debug("participant connected", slog.String("participant_id", participant.SID()), slog.String("identity", participant.Identity()))
			},
			OnParticipantDisconnected: func(participant *lksdk.RemoteParticipant) {
				c.config.Logger.Debug("participant disconnected", slog.String("participant_id", participant.SID()), slog.String("identity", participant.Identity()))
			},
			OnActiveSpeakersChanged: func(participants []lksdk.Participant) {
				c.config.Logger.Debug("active speakers changed", slog.Int("count", len(participants)))
			},
			OnRoomMetadataChanged: func(metadata string) {
				c.config.Logger.Debug("room metadata changed", slog.String("metadata", metadata))
			},
			OnRecordingStatusChanged: func(isRecording bool) {
				c.config.Logger.Debug("recording status changed", slog.Bool("is_recording", isRecording))
			},
			OnRoomMoved: func(roomName string, token string) {
				c.config.Logger.Debug("room moved", slog.String("room_name", roomName))
			},
			OnReconnecting: func() {
				c.config.Logger.Debug("reconnecting to voice")
			},
			OnReconnected: func() {
				c.config.Logger.Debug("reconnected to voice")
			},
			OnLocalTrackSubscribed: func(publication *lksdk.LocalTrackPublication, lp *lksdk.LocalParticipant) {
				c.config.Logger.Debug("local track subscribed", slog.String("track_sid", publication.SID()), slog.String("participant_id", lp.SID()))
			},
			ParticipantCallback: lksdk.ParticipantCallback{
				OnLocalTrackPublished: func(publication *lksdk.LocalTrackPublication, lp *lksdk.LocalParticipant) {
					c.config.Logger.Debug("local track published", slog.String("track_sid", publication.SID()), slog.String("participant_id", lp.SID()))
				},
				OnLocalTrackUnpublished: func(publication *lksdk.LocalTrackPublication, lp *lksdk.LocalParticipant) {
					c.config.Logger.Debug("local track unpublished", slog.String("track_sid", publication.SID()), slog.String("participant_id", lp.SID()))
				},
				OnTrackMuted: func(pub lksdk.TrackPublication, p lksdk.Participant) {
					c.config.Logger.Debug("track muted", slog.String("track_sid", pub.SID()), slog.String("participant_id", p.SID()))
				},
				OnTrackUnmuted: func(pub lksdk.TrackPublication, p lksdk.Participant) {
					c.config.Logger.Debug("track unmuted", slog.String("track_sid", pub.SID()), slog.String("participant_id", p.SID()))
				},
				OnMetadataChanged: func(oldMetadata string, p lksdk.Participant) {
					c.config.Logger.Debug("metadata changed", slog.String("participant_id", p.SID()), slog.String("old_metadata", oldMetadata))
				},
				OnAttributesChanged: func(changed map[string]string, p lksdk.Participant) {
					c.config.Logger.Debug("attributes changed", slog.String("participant_id", p.SID()), slog.Any("changed", changed))
				},
				OnIsSpeakingChanged: func(p lksdk.Participant) {
					c.config.Logger.Debug("is speaking changed", slog.String("participant_id", p.SID()), slog.Bool("is_speaking", p.IsSpeaking()))
				},
				OnConnectionQualityChanged: func(update *livekit.ConnectionQualityInfo, p lksdk.Participant) {
					c.config.Logger.Debug("connection quality changed", slog.String("participant_id", p.SID()), slog.Any("quality", update.Quality))
				},
				OnTrackSubscribed: func(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
					c.config.Logger.Debug("track subscribed", slog.String("track_sid", publication.SID()), slog.String("participant_id", rp.SID()))
				},
				OnTrackUnsubscribed: func(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
					c.config.Logger.Debug("track unsubscribed", slog.String("track_sid", publication.SID()), slog.String("participant_id", rp.SID()))
				},
				OnTrackSubscriptionFailed: func(sid string, rp *lksdk.RemoteParticipant) {
					c.config.Logger.Debug("track subscription failed", slog.String("track_sid", sid), slog.String("participant_id", rp.SID()))
				},
				OnTrackPublished: func(publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
					c.config.Logger.Debug("track published", slog.String("track_sid", publication.SID()), slog.String("participant_id", rp.SID()))
				},
				OnTrackUnpublished: func(publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
					c.config.Logger.Debug("track unpublished", slog.String("track_sid", publication.SID()), slog.String("participant_id", rp.SID()))
				},
				OnDataReceived: func(data []byte, params lksdk.DataReceiveParams) {
					c.config.Logger.Debug("data received", slog.Int("size", len(data)), slog.String("participant_id", params.SenderIdentity))
				},
				OnDataPacket: func(data lksdk.DataPacket, params lksdk.DataReceiveParams) {
					c.config.Logger.Debug("data packet received", slog.String("participant_id", params.SenderIdentity))
				},
				OnTranscriptionReceived: func(transcriptionSegments []*lksdk.TranscriptionSegment, p lksdk.Participant, publication lksdk.TrackPublication) {
					c.config.Logger.Debug("transcription received", slog.String("participant_id", p.SID()), slog.Int("segments", len(transcriptionSegments)))
				},
			},
		})

		if err := room.JoinWithToken(c.state.Endpoint, c.state.Token); err != nil {
			c.config.Logger.Error("error connecting to voice", slog.Any("err", err))
			return
		}

		c.room = room
		c.opened()
	}
	c.state.SessionID = update.SessionID
}

func (c *connImpl) HandleVoiceServerUpdate(update gateway.EventVoiceServerUpdate) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if update.GuildID != c.state.GuildID || update.Endpoint == nil {
		return
	}

	c.state.Token = update.Token
	c.state.ConnectionID = update.ConnectionID
	c.state.Endpoint = *update.Endpoint
	c.state.ChannelID = &update.ChannelID

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := c.voiceStateUpdateFunc(ctx, gateway.MessageDataVoiceStateUpdate{
			GuildID:      c.state.GuildID,
			ChannelID:    c.state.ChannelID,
			ConnectionID: &c.state.ConnectionID,
		}); err != nil {
			c.config.Logger.Error("error sending voice state update to connect to voice gateway", slog.Any("err", err))
		}
	}()
}

func (c *connImpl) Open(ctx context.Context, channelID snowflake.ID, selfMute bool, selfDeaf bool) error {
	c.config.Logger.Debug("opening voice conn")

	if err := c.voiceStateUpdateFunc(ctx, gateway.MessageDataVoiceStateUpdate{
		GuildID:      c.state.GuildID,
		ChannelID:    &channelID,
		SelfMute:     selfMute,
		SelfDeaf:     selfDeaf,
		SelfVideo:    false,
		SelfStream:   false,
		ConnectionID: nil,
	}); err != nil {
		return err
	}

	select {
	case <-c.openedCtx.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *connImpl) Close(ctx context.Context) {
	if err := c.voiceStateUpdateFunc(ctx, gateway.MessageDataVoiceStateUpdate{
		GuildID: c.state.GuildID,
	}); err != nil {
		c.config.Logger.Error("error sending voice state update to close voice conn", slog.Any("err", err))
	}
	defer c.room.Disconnect()

	select {
	case <-c.closedCtx.Done():
	case <-ctx.Done():
	}
	c.removeConnFunc()
}

func (c *connImpl) Room() *lksdk.Room {
	return c.room
}
