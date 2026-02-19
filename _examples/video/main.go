package main

import "C"
import (
	"bufio"
	"context"
	_ "embed"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/jonas747/ogg"

	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/voice"
)

var (
	token            = os.Getenv("fluxergo_token")
	guildID          = snowflake.GetEnv("fluxergo_guild_id")
	channelID        = snowflake.GetEnv("fluxergo_channel_id")
	messageChannelID = snowflake.GetEnv("fluxergo_message_channel_id")
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("starting up")
	slog.Info("fluxergo version", slog.String("version", fluxergo.Version))

	client, err := fluxergo.New(token,
		bot.WithDefaultGateway(),
		// bot.WithEventListenerFunc(func(e *events.Ready) {
		// 	go play(e.Client(), guildID, channelID, messageChannelID, "https://www.youtube.com/watch?v=wJNbtYdr-Hg")
		// }),
		bot.WithEventListenerFunc(onMessage),
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

	if err = client.OpenGateway(context.Background()); err != nil {
		slog.Error("error connecting to gateway", slog.Any("error", err))
		return
	}

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func onMessage(e *events.MessageCreate) {
	content := e.Message.Content
	if !strings.HasPrefix(content, "!") {
		return
	}

	content = strings.TrimPrefix(content, "!")
	fields := strings.FieldsFunc(content, func(r rune) bool {
		return r == ' '
	})

	command := fields[0]
	args := fields[1:]

	switch command {
	case "ping":
		e.Client().Rest.CreateMessage(e.Message.ChannelID, fluxer.MessageCreate{
			Content: "Pong!",
		})
	case "play":
		if len(args) == 0 {
			e.Client().Rest.CreateMessage(e.Message.ChannelID, fluxer.MessageCreate{
				Content: "Please provide a URL to play.",
			})
			return
		}
		url := args[0]
		go play(e.Client(), *e.GuildID, 1471575381777371381, e.ChannelID, url)
	case "stop":
		conn := e.Client().VoiceManager.GetConn(*e.GuildID)
		if conn != nil {
			conn.Close(context.Background())
		}
	}
}

func play(client *bot.Client, guildID snowflake.ID, channelID snowflake.ID, messageChannelID snowflake.ID, url string) {
	conn := client.VoiceManager.CreateConn(guildID)

	videoReader, audioReader, err := startStream(url)
	if err != nil {
		client.Rest.CreateMessage(messageChannelID, fluxer.MessageCreate{
			Content: "Error starting stream: " + err.Error(),
		})
		return
	}

	if err = conn.Open(context.Background(), channelID); err != nil {
		client.Rest.CreateMessage(messageChannelID, fluxer.MessageCreate{
			Content: "Error connecting to voice channel: " + err.Error(),
		})
		return
	}

	video, err := conn.LiveKit().VideoWriter("video", voice.VideoSourceScreenShare, 1280, 720, 30)
	if err != nil {
		client.Rest.CreateMessage(messageChannelID, fluxer.MessageCreate{
			Content: "Error creating video writer: " + err.Error(),
		})
	}
	defer video.Close()

	audio, err := conn.LiveKit().AudioWriter("audio", voice.AudioSourceScreenShare)
	if err != nil {
		client.Rest.CreateMessage(messageChannelID, fluxer.MessageCreate{
			Content: "Error creating audio writer: " + err.Error(),
		})
	}
	defer audio.Close()

	go writeVideo(video, videoReader)
	go writeAudio(audio, audioReader)

	select {}
}

func startStream(url string) (video io.Reader, audio io.Reader, err error) {
	// yt-dlp -> stdout
	ytdlp := exec.Command(
		"yt-dlp",
		"-f", "bv*+ba/b",
		"--quiet",
		"--no-progress",
		"-o", "-",
		url,
	)

	ytdlpOut, err := ytdlp.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err = ytdlp.Start(); err != nil {
		return nil, nil, err
	}

	// Create extra pipe for ffmpeg audio output
	audioPipeReader, audioPipeWriter, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	ffmpeg := exec.Command(
		"ffmpeg",
		"-re",
		"-i", "pipe:0",

		// Video output
		"-map", "0:v:0",
		"-c:v", "libx264",
		"-profile:v", "baseline",
		"-pix_fmt", "yuv420p",
		"-preset", "veryfast",
		"-tune", "zerolatency",
		"-g", "60",
		"-keyint_min", "60",
		"-x264opts", "no-scenecut:repeat-headers=1",
		"-b:v", "2500k",
		"-maxrate", "2500k",
		"-bufsize", "5000k",
		"-f", "h264",
		"pipe:1",

		// Audio output (to fd 3)
		"-map", "0:a:0",
		"-c:a", "libopus",
		"-ac", "2",
		"-ar", "48000",
		"-b:a", "128K",
		"-f", "ogg",
		"pipe:3",
	)

	ffmpeg.Stdin = ytdlpOut
	ffmpeg.ExtraFiles = []*os.File{audioPipeWriter}

	videoPipe, err := ffmpeg.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err = ffmpeg.Start(); err != nil {
		return nil, nil, err
	}

	// Close writer in parent so ffmpeg owns it
	audioPipeWriter.Close()

	return videoPipe, audioPipeReader, nil
}

func writeVideo(w io.Writer, r io.Reader) {
	reader := bufio.NewReader(r)

	var nalBuf []byte
	startCode := []byte{0x00, 0x00, 0x00, 0x01}

	for {
		chunk := make([]byte, 4096)
		n, err := reader.Read(chunk)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read error:", err)
			continue
		}

		nalBuf = append(nalBuf, chunk[:n]...)

		for {
			idx := indexOf(nalBuf[4:], startCode)
			if idx == -1 {
				break
			}

			// safe copy of one NAL frame
			frame := make([]byte, idx+4)
			copy(frame, nalBuf[:idx+4])
			nalBuf = nalBuf[idx+4:]

			if _, err := w.Write(frame); err != nil {
				log.Println("write frame error:", err)
				return
			}
		}
	}
}

func writeAudio(w io.Writer, r io.Reader) {
	decoder := ogg.NewPacketDecoder(ogg.NewDecoder(bufio.NewReader(r)))

	for {
		data, _, err := decoder.Decode()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("decode error:", err)
			continue
		}

		if _, err = w.Write(data); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}

func indexOf(data, sep []byte) int {
	for i := 0; i <= len(data)-len(sep); i++ {
		if string(data[i:i+len(sep)]) == string(sep) {
			return i
		}
	}
	return -1
}
