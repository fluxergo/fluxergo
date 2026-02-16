package voice

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/gateway"
)

type (
	// ConnCreateFunc is a type alias for a function that creates a new Conn.
	ConnCreateFunc func(guildID snowflake.ID, userID snowflake.ID, voiceStateUpdateFunc StateUpdateFunc, removeConnFunc func(), opts ...ConnConfigOpt) Conn

	// Conn is a complete voice conn to fluxer. It holds the Gateway and voiceudp.UDPConn conn and combines them.
	Conn interface {
		LiveKit() *LivekitConn

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
	}
)

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

	conn.liveKitConn = cfg.LiveKitConnCreateFunc()

	return conn
}

type connImpl struct {
	config               connConfig
	voiceStateUpdateFunc StateUpdateFunc
	removeConnFunc       func()

	state   State
	stateMu sync.Mutex

	liveKitConn *LivekitConn

	audioSender   AudioSender
	audioReceiver AudioReceiver

	openedCtx context.Context
	opened    context.CancelFunc
	closedCtx context.Context
	closed    context.CancelFunc
}

func (c *connImpl) LiveKit() *LivekitConn {
	return c.liveKitConn
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
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if update.GuildID != c.state.GuildID || update.UserID != c.state.UserID {
		return
	}

	c.state.SessionID = update.SessionID
	if update.ChannelID == nil {
		c.state.ChannelID = nil
		c.liveKitConn.Close()
		c.closed()
	} else {
		c.state.ChannelID = update.ChannelID

		if err := c.liveKitConn.Open(c.state); err != nil {
			c.config.Logger.Error("error connecting to voice", slog.Any("err", err))
			return
		}
		c.opened()
	}
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
	defer c.liveKitConn.Close()

	select {
	case <-c.closedCtx.Done():
	case <-ctx.Done():
	}
	c.removeConnFunc()
}
