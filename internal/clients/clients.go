package clients

import (
	"url-shortener/internal/clients/kafka"
	ssogrpc "url-shortener/internal/clients/sso/grpc"
)

type Clients struct {
	*ssogrpc.Client
	*kafka.Broker
	appID int32
}

func New(sso *ssogrpc.Client, kafka *kafka.Broker) Clients {
	return Clients{
		Client: sso,
		Broker: kafka,
		appID:  sso.AppID,
	}
}

func (c Clients) AppID() int32 {
	return c.appID
}
