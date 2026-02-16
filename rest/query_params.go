package rest

import (
	"github.com/fluxergo/fluxergo/fluxer"
)

// QueryParams serves as a generic interface for implementations of rest endpoint query parameters.
type QueryParams interface {
	// ToQueryValues transforms fields from the QueryParams interface implementations into fluxer.QueryValues.
	ToQueryValues() fluxer.QueryValues
}
