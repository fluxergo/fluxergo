package main

import "C"
import (
	"context"
	_ "embed"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/voice"
)

var (
	token     = os.Getenv("fluxergo_token")
	guildID   = snowflake.GetEnv("fluxergo_guild_id")
	channelID = snowflake.GetEnv("fluxergo_channel_id")
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
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

	time.Sleep(2 * time.Second)

	videoTrack, err := conn.LiveKit().VideoWriter("video", voice.VideoSourceScreenShare)
	if err != nil {
		panic("error creating video writer: " + err.Error())
	}
	defer videoTrack.Close()
	slog.Info("created video writer")

	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-i", "_examples/video/test.mp4",
		"-map", "0:v:0",
		"-c:v", "libx264",
		"-profile:v", "baseline",
		"-pix_fmt", "yuv420p",
		"-tune", "zerolatency",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"pipe:1",
	)

	videoPipe, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal("error starting ffmpeg:", err)
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := videoPipe.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("video read error:", err)
			continue
		}

		chunk := make([]byte, n)
		copy(chunk, buf[:n])

		videoTrack.Write(chunk)
	}

	cmd.Wait()
	select {}
}
