package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
	"github.com/fluxergo/fluxergo/events"
	"github.com/fluxergo/fluxergo/fluxer"
	"github.com/fluxergo/fluxergo/gateway"
	"github.com/fluxergo/fluxergo/rest"
)

var token = os.Getenv("fluxergo_token")

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.Info("starting up")
	slog.Info("fluxergo version", slog.String("version", fluxergo.Version))

	client, err := fluxergo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithURL("wss://gateway.fluxer.app"),
		),
		bot.WithRestClientConfigOpts(
			rest.WithURL("https://api.fluxer.app/v1"),
		),
		bot.WithEventListenerFunc(onMessage),
	)
	if err != nil {
		slog.Error("error creating client", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error connecting to gateway", slog.Any("error", err))
		return
	}

	slog.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func onMessage(e *events.MessageCreate) {
	if e.Message.Content == "!ping" {
		if _, err := e.Client().Rest.CreateMessage(e.ChannelID, fluxer.MessageCreate{
			Content: "Pong!",
		}); err != nil {
			slog.Error("error sending message", slog.Any("err", err))
		}
	}
}
