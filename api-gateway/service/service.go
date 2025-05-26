package service

import (
	"api-gateway/client/consul"
	"api-gateway/client/http_adapter"
)

type Service struct {
	httpClient      *http_adapter.Client
	discoveryClient *consul.DiscoveryClient
}

func NewService(httpClient *http_adapter.Client, discoveryClient *consul.DiscoveryClient) *Service {
	return &Service{
		httpClient:      httpClient,
		discoveryClient: discoveryClient,
	}
}
