package events

import "github.com/fluxergo/fluxergo/gateway"

type HeartbeatAck struct {
	*GenericEvent
	gateway.EventHeartbeatAck
}
