package main

import "C"
import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/livekit/protocol/logger"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"

	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
)

var (
	token     = os.Getenv("fluxergo_token")
	guildID   = snowflake.GetEnv("fluxergo_guild_id")
	channelID = snowflake.GetEnv("fluxergo_channel_id")

	//go:embed nico.dca
	testOpus []byte
)

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("starting up")
	slog.Info("fluxergo version", slog.String("version", fluxergo.Version))

	client, err := fluxergo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			go play(e.Client())
		}),
	)
	if err != nil {
		slog.Error("error creating client", slog.Any("err", err))
		return
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		client.Close(ctx)
	}()

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error connecting to gateway", slog.Any("error", err))
		return
	}

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func play(client *bot.Client) {
	conn := client.VoiceManager.CreateConn(guildID)
	slog.Info("connecting to voice manager")

	if err := conn.Open(context.Background(), channelID, false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}
	slog.Info("connected to voice channel")

	room := conn.Room()
	for {
		if room.ConnectionState() != lksdk.ConnectionStateConnected {
			slog.Warn("connection state is ", room.ConnectionState())
			time.Sleep(time.Second)
		} else {
			slog.Info("CONNECTED")
			break
		}
	}

	slog.Info("creating track")

	track, err := lksdk.NewLocalTrack(webrtc.RTPCodecCapability{
		MimeType:    webrtc.MimeTypeOpus,
		ClockRate:   48000,
		Channels:    2,
		SDPFmtpLine: "",
	})
	if err != nil {
		panic("error creating track: " + err.Error())
	}

	slog.Info("publishing track")
	pub, err := room.LocalParticipant.PublishTrack(track, &lksdk.TrackPublicationOptions{
		Name: "lol",
	})
	if err != nil {
		panic("error publishing track: " + err.Error())
	}

	slog.Info("published track", slog.String("track_id", pub.SID()))

	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var frameLen int16
	for {
		r := bytes.NewReader(testOpus)
		for range ticker.C {
			err = binary.Read(r, binary.LittleEndian, &frameLen)
			if err != nil {
				if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
					break
				}
				panic("error reading file: " + err.Error())
			}
			if frameLen <= 0 {
				panic("invalid frame length: " + strconv.Itoa(int(frameLen)))
			}

			buf := make([]byte, frameLen)
			_, err = r.Read(buf)
			if err != nil {
				break
			}

			err = track.WriteSample(media.Sample{
				Data:      buf,
				Duration:  20 * time.Millisecond,
				Timestamp: time.Now(),
			}, nil)
			if err != nil {
				logger.Errorw("error writing sample", err)
			}
		}
	}

}
