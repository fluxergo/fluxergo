package events

import (
	"github.com/fluxergo/fluxergo/fluxer"
)

// SelfUpdate is called when something about this fluxer.User updates
type SelfUpdate struct {
	*GenericEvent
	SelfUser    fluxer.OAuth2User
	OldSelfUser fluxer.OAuth2User
}
