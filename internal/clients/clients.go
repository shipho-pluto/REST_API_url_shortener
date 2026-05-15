package clients

import ssogrpc "url-shortener/internal/clients/sso/grpc"

type Clients struct {
	sso *ssogrpc.Client
}

func New(sso *ssogrpc.Client) Clients {
	return Clients{
		sso: sso,
	}
}
