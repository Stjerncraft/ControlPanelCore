package core

import (
	"github.com/stjerncraft/controlpanelcore/api"
)

type ServiceProvider struct {
	name         string
	providedApis []*api.Api
	server       *MinecraftServer //The Minecraft Server this belongs to. Can be nil for global Service Providers

	core    *Core
	id      int32
	handler ServiceProviderHandler
}

func newServiceProvider(name string, handler ServiceProviderHandler, c *Core, server *MinecraftServer) *ServiceProvider {
	newProvider := new(ServiceProvider)
	newProvider.name = name
	newProvider.handler = handler
	newProvider.core = c
	newProvider.server = server

	return newProvider
}

func (s *ServiceProvider) GetApiList() *[]*api.Api {
	return &s.providedApis
}

func (s *ServiceProvider) GetName() string {
	return s.name
}

func (s *ServiceProvider) GetHandler() ServiceProviderHandler {
	return s.handler
}

func (s *ServiceProvider) GetCore() *Core {
	return s.core
}

func (s *ServiceProvider) GetServer() *MinecraftServer {
	return s.server
}

func (s *ServiceProvider) GetId() int32 {
	return s.id
}
