package core

import (
	"errors"
	"github.com/stjerncraft/controlpanelcore/api"
	"github.com/stjerncraft/controlpanelcore/config"
)

type Core struct {
	Config config.Config

	ApiMap        map[string]*api.Api //Json Api to Api
	ApiServiceMap map[*api.Api][]*ServiceProvider

	nextProviderId      int32 //Unique id for ServiceProvider
	ServiceProviderList map[int32]*ServiceProvider

	MinecraftServers map[string]*MinecraftServer
	//Events:
	//EventNewServiceProvider
	//EventNewServiceProviderForApi
	//EventNewMinecraftProvider
}

func NewCore(config config.Config) (*Core, error) {
	newCore := new(Core)
	newCore.Config = config

	newCore.ApiMap = make(map[string]*api.Api)
	newCore.ApiServiceMap = make(map[*api.Api][]*ServiceProvider)

	newCore.nextProviderId = 0
	newCore.ServiceProviderList = make(map[int32]*ServiceProvider)

	newCore.MinecraftServers = make(map[string]*MinecraftServer)

	newCore.AddServiceProvider("CoreServiceProvider", NewCoreServiceProviderHandler(), nil)

	return newCore, nil
}

func (c *Core) AddServiceProvider(name string, serviceHandler ServiceProviderHandler, server *MinecraftServer) *ServiceProvider {
	newServiceProvider := newServiceProvider(name, serviceHandler, c, server)

	providerId := c.nextProviderId
	newServiceProvider.id = providerId
	c.nextProviderId++

	c.ServiceProviderList[providerId] = newServiceProvider
	serviceHandler.OnAdd(newServiceProvider)

	return newServiceProvider
}

func (c *Core) RegisterApiServiceProvider(api *api.Api, serviceIn *ServiceProvider) error {
	//Register the API if it doesn't already exist
	actualApi := c.addApi(api)
	if actualApi != api {
		return errors.New("Assert failure: Service Provider has different Api instance!")
	}

	//Service can only provide for one Api once
	serviceList, gotList := c.ApiServiceMap[api]
	if !gotList {
		serviceList = make([]*ServiceProvider, 1)
		serviceList[0] = serviceIn
		c.ApiServiceMap[api] = serviceList
	} else {
		c.ApiServiceMap[api] = append(serviceList, serviceIn)
	}
	for _, curService := range serviceList {
		if curService == serviceIn {
			//TODO: Just silently fail instead?
			return errors.New("Assert failure: " + serviceIn.GetName() + " already registered for api: " + *api.GetJson())
		}
	}
	serviceIn.GetHandler().OnAddApi(api)

	return nil
}

func (c *Core) GetApi(json *string) *api.Api {
	return c.ApiMap[*json]
}

//Add api, returning existing Api if already defined
func (c *Core) addApi(api *api.Api) *api.Api {
	existingApi := c.GetApi(api.GetJson())
	if existingApi == nil {
		c.ApiMap[*api.GetJson()] = api
		existingApi = api
	}
	return existingApi
}

func (c *Core) AddMinecraftServer(server *MinecraftServer) error {
	if _, ok := c.MinecraftServers[server.Id]; ok {
		return errors.New("Minecraft Server already exists with the given id: " + server.Id)
	}

	//Check if the responsible Service implementes the correct API for control

	c.MinecraftServers[server.Id] = server

	//TODO: Send out Event for new server

	return nil
}
