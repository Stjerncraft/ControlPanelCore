package core

import (
	"github.com/stjerncraft/controlpanelcore/api"
)

type ServiceProviderHandler interface {

	//Added as handler
	OnAdd(serviceProvider *ServiceProvider)
	//Removed as handler
	OnRemove()

	//Added a handler for api
	OnAddApi(api *api.Api)
	//Removed as handler for api
	OnRemoveApi(api *api.Api)

	//A Module has started a new session
	OnSessionStart()
	//A Module has ended it's session
	OnSessionEnd()

	//A session registered as listener for event
	OnEventRegister() error
	//A session unregistered from event
	OnEventUnregister()

	//A method call for a given session
	OnMethodCall(method *api.ApiMethod, args *[]api.ApiValue) api.ApiValue

	//HTTP Request for binary resource
	OnResourceRequest(name string, args []string)
}
