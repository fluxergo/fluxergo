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

	videoReader, err := getReader(url, "bv*")
	if err != nil {
		client.Rest.CreateMessage(messageChannelID, fluxer.MessageCreate{
			Content: "Error getting video reader: " + err.Error(),
		})
	}

	audioReader, err := getReader(url, "ba")
	if err != nil {
		client.Rest.CreateMessage(messageChannelID, fluxer.MessageCreate{
			Content: "Error getting audio reader: " + err.Error(),
		})
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

func getReader(url string, format string) (io.Reader, error) {
	cmd := exec.Command(
		"yt-dlp",
		"-f", format,
		"--quiet",
		"--no-progress",
		"-o", "-",
		url,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		data, _ := io.ReadAll(stdErr)
		if len(data) > 0 {
			slog.Error("yt-dlp stderr", slog.String("format", format), slog.String("output", string(data)))
		}
		if err = cmd.Wait(); err != nil {
			slog.Error("error waiting for yt-dlp", slog.String("format", format), slog.Any("err", err))
		}
	}()

	return stdout, nil
}

func writeVideo(w io.Writer, r io.Reader) {
	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-i", "pipe:0",
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
		"-an",
		"-f", "h264",
		"pipe:1",
	)

	cmd.Stdin = r

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("error creating stdout pipe", slog.Any("err", err))
		return
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("error creating stderr pipe", slog.Any("err", err))
		return
	}

	if err = cmd.Start(); err != nil {
		slog.Error("error starting command", slog.Any("err", err))
		return
	}
	go func() {
		data, _ := io.ReadAll(stdErr)
		if len(data) > 0 {
			slog.Error("ffmpeg stderr", slog.String("output", string(data)))
		}
		if err = cmd.Wait(); err != nil {
			slog.Error("failed to wait video ffmpeg", slog.Any("error", err))
		}
	}()

	reader := bufio.NewReader(pipe)
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
	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-i", "pipe:0",
		"-c:a", "libopus",
		"-ac", "2",
		"-ar", "48000",
		"-b:a", "64K",
		"-f", "ogg",
		"-vn",
		"pipe:1",
	)
	cmd.Stdin = r

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("error creating stdout pipe", slog.Any("err", err))
		return
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("error creating stderr pipe", slog.Any("err", err))
		return
	}

	decoder := ogg.NewPacketDecoder(ogg.NewDecoder(bufio.NewReaderSize(pipe, 65307)))

	if err = cmd.Start(); err != nil {
		slog.Error("error starting video", slog.Any("err", err))
		return
	}
	go func() {
		data, _ := io.ReadAll(stdErr)
		if len(data) > 0 {
			slog.Error("ffmpeg stderr", slog.String("output", string(data)))
		}
		if err = cmd.Wait(); err != nil {
			slog.Error("failed to wait audio ffmpeg", slog.Any("error", err))
		}
	}()

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
