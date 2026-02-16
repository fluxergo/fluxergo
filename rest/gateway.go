package rest

import (
	"github.com/fluxergo/fluxergo/fluxer"
)

var _ Gateway = (*gatewayImpl)(nil)

func NewGateway(client Client) Gateway {
	return &gatewayImpl{client: client}
}

type Gateway interface {
	GetGatewayBot(opts ...RequestOpt) (*fluxer.GatewayBot, error)
}

type gatewayImpl struct {
	client Client
}

func (s *gatewayImpl) GetGatewayBot(opts ...RequestOpt) (gatewayBot *fluxer.GatewayBot, err error) {
	err = s.client.Do(GetGatewayBot.Compile(nil), nil, &gatewayBot, opts...)
	return
}
