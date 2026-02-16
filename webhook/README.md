# webhook

[Webhook](https://fluxer.com/developers/docs/resources/webhook) module of [disgo](https://github.com/fluxergo/fluxergo)

### Usage

Import the package into your project.

```go
import "github.com/fluxergo/fluxergo/webhook"
```

Create a new Webhook by `webhook_id` and `webhook_token`. (*This WebhookClient should be created once as it holds important state*)

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken")

client, err := webhook.NewWithURL("webhookURL")
```

`webhook.New` takes a vararg of type `webhook.ConfigOpt` as third argument which lets you pass additional optional parameter like a custom logger, rest client, etc

### Optional Arguments

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken",
	webhook.WithLogger(logrus.New()),
	webhook.WithDefaultAllowedMentions(fluxer.AllowedMentions{
		RepliedUser: false,
	}),
)
```

### Send Message

You can send a message as following

```go
client := webhook.New(snowflake.ID("webhookID"), "webhookToken")

message, err := client.CreateContent("hello world!")

message, err := client.CreateEmbeds(fluxer.NewEmbedBuilder().
	SetDescription("hello world!").
	Build(),
)

message, err := client.CreateMessage(fluxer.NewWebhookMessageCreateBuilder().
	SetContent("hello world!").
	Build(),
	rest.CreateWebhookMessageParams{},
)

message, err := client.CreateMessage(fluxer.WebhookMessageCreate{
	Content: "hello world!",
}, rest.CreateWebhookMessageParams{})
```

### Full Example

a full example can be found [here](../_examples/webhook/main.go)