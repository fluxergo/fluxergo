package events

import "github.com/fluxergo/fluxergo/gateway"

// Ready indicates we received the Ready from the gateway.Gateway
type Ready struct {
	*GenericEvent
	gateway.EventReady
}

// Resumed indicates fluxergo resumed the gateway.Gateway
type Resumed struct {
	*GenericEvent
}
