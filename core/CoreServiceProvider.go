package core

import (
	"fmt"
	"github.com/stjerncraft/controlpanelcore/api"
	"github.com/stjerncraft/controlpanelcore/json"
)

var CoreServiceApiJson = `{
	"name": "CoreAPI v1.0.0",

	"objects": {
		"MinecraftServer": {
			"Id": "string",
			"Name": "string",
			"Address": "string",
			"Status": "string",
			"PlayersOnline": "int",
			"PlayersMax": "int",
			"ProviderId": "int"
		}
	},
	"events": {
		"AddedServer": ["MinecraftServer"],
		"RemovedServer": ["MinecraftServer"],
		"UpdatedServer": ["MinecraftServer"]
	},
	"methods": {
		"addServer": {
			"args": ["MinecraftServer"],
			"return": "int"
		},
		"updateServer": {
			"args": ["MinecraftServer"],
			"return": "byte"
		},
		"removeServer": {
			"args": ["int"],
			"return": "byte"
		},
		"getServer": {
			"args": ["int"],
			"return": "MinecraftServer"
		},
		"getAllServers": {
			"args": [],
			"return": "list[MinecraftServer]"
		}
	}
}`

type CoreServiceProviderHandler struct {
	serviceProvider *ServiceProvider
	api             *api.Api
}

func NewCoreServiceProviderHandler() *CoreServiceProviderHandler {
	return new(CoreServiceProviderHandler)
}

func (c *CoreServiceProviderHandler) OnAdd(serviceProvider *ServiceProvider) {
	c.serviceProvider = serviceProvider

	//Register to API
	jsonApi, err := json.ParseJsonApi(CoreServiceApiJson)
	if err != nil {
		//TODO: Remove the service on failure
		//serviceProvider.core.RemoveService()
		fmt.Println("Failed to create API for service: ", err)
		return
	}

	newApi, err := api.ApiFromJson(jsonApi)
	if err != nil {
		//TODO: Remove service
		fmt.Println("Failed to create API from json API: ", err)
		return
	}

	c.api = newApi
}

func (c *CoreServiceProviderHandler) OnRemove() {

}

func (c *CoreServiceProviderHandler) OnAddApi(api *api.Api) {
	if api != c.api {
		//TODO: Unregister from API
		fmt.Printf("Incorrect api for CoreServiceProviderHandler: %+v\n", api)
		return
	}
}

func (c *CoreServiceProviderHandler) OnRemoveApi(api *api.Api) {

}

func (c *CoreServiceProviderHandler) OnSessionStart() {

}

func (c *CoreServiceProviderHandler) OnSessionEnd() {

}

func (c *CoreServiceProviderHandler) OnEventRegister() error {
	return nil
}

func (c *CoreServiceProviderHandler) OnEventUnregister() {

}

func (c *CoreServiceProviderHandler) OnMethodCall(method *api.ApiMethod, args *[]api.ApiValue) api.ApiValue {
	switch method.Name {
	case "addServer":
		server := toMinecraftServer((*args)[0].(*api.ApiValueObject))
		c.serviceProvider.core.AddMinecraftServer(server)
	}

	return nil
}

func (c *CoreServiceProviderHandler) OnResourceRequest(name string, args []string) {

}

func toMinecraftServer(obj *api.ApiValueObject) *MinecraftServer {
	id := obj.Fields["Id"].(*api.ApiValueString)
	name := obj.Fields["Name"].(*api.ApiValueString)
	address := obj.Fields["Address"].(*api.ApiValueString)
	status := obj.Fields["Status"].(*api.ApiValueString)
	playersOnline := obj.Fields["PlayersOnline"].(*api.ApiValueInteger)
	playersMax := obj.Fields["PlayersMax"].(*api.ApiValueInteger)
	providerId := obj.Fields["ProviderId"].(*api.ApiValueInteger)

	return &MinecraftServer{
		Id:            id.Value,
		Name:          name.Value,
		Address:       address.Value,
		Status:        status.Value,
		PlayersOnline: playersOnline.Value,
		PlayersMax:    playersMax.Value,
		ProviderId:    providerId.Value,
	}
}
