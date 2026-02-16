package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
)

var (
	token = os.Getenv("fluxergo_token")
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(logger))
	slog.Info("starting example...")
	slog.Info("FluxerGo version", slog.Any("version", fluxergo.Version))

	client, err := fluxergo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithPresenceOpts(gateway.WithListeningActivity("your bullshit", gateway.WithActivityState("lol")), gateway.WithOnlineStatus(fluxer.OnlineStatusDND)),
		),
		bot.WithEventListenerFunc(onMessage),
	)
	if err != nil {
		slog.Error("error while building bot instance", slog.Any("err", err))
		return
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to discord", slog.Any("err", err))
	}

	defer client.Close(context.TODO())

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessage(e *events.MessageCreate) {
	if e.Message.Content != "!ping" {
		return
	}
	if _, err := e.Client().Rest.CreateMessage(e.ChannelID, fluxer.NewMessageCreate().WithContent("pong")); err != nil {
		slog.Error("error while sending message", slog.Any("err", err))
	}
}
