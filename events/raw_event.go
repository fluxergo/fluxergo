package events

import "github.com/fluxergo/fluxergo/gateway"

type Raw struct {
	*GenericEvent
	gateway.EventRaw
}
