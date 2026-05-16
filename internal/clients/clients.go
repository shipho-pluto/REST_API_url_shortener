package clients

import ssogrpc "url-shortener/internal/clients/sso/grpc"

type Clients struct {
	*ssogrpc.Client
	appID int32
}

func New(sso *ssogrpc.Client, appID int32) Clients {
	return Clients{
		Client: sso,
		appID:  appID,
	}
}

func (c Clients) AppID() int32 {
	return c.appID
}
