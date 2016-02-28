package api

import "github.com/stjerncraft/controlpanelcore/json"

type Api struct {
	name    string
	apiJson *string //Used for comparison and lookup

	objects map[string]*ApiObject
	events  map[string]*ApiEvent
	methods map[string]*ApiMethod
}

func ApiFromJson(jsonApi *json.JsonApi) (*Api, error) {
	newApi := new(Api)
	newApi.name = jsonApi.Name

	apiJson, err := json.GetJsonApi(jsonApi)
	if err != nil {
		return nil, err
	}
	newApi.apiJson = &apiJson

	//Create a list of Objects, allowing them to reference each other for later stages
	newApi.objects = make(map[string]*ApiObject)
	for name, _ := range jsonApi.Objects {
		newApi.objects[name] = NewApiObject(name, newApi)
	}
	//Populate fields
	for name, obj := range jsonApi.Objects {
		err := newApi.objects[name].SetFields((*map[string]string)(&obj))
		if err != nil {
			return nil, err
		}
	}
	//Inherit(All Objects now have their fields populated)
	for _, obj := range newApi.objects {
		err := obj.InheritFields()
		if err != nil {
			return nil, err
		}
	}

	//Events
	newApi.events = make(map[string]*ApiEvent)
	for name, ev := range jsonApi.Events {
		event, err := NewApiEvent((*[]string)(&ev), newApi)
		if err != nil {
			return nil, err
		}
		newApi.events[name] = event
	}

	//Methods
	newApi.methods = make(map[string]*ApiMethod)
	for name, method := range jsonApi.Methods {
		newMethod, err := NewApiMethod(name, method.Args, method.Return, newApi)
		if err != nil {
			return nil, err
		}

		newApi.methods[name] = newMethod
	}

	return newApi, nil
}

func (api *Api) GetName() string {
	return api.name
}

func (api *Api) GetObjects() *map[string]*ApiObject {
	return &api.objects
}

func (api *Api) GetObject(name string) *ApiObject {
	return api.objects[name]
}

func (api *Api) GetMethods() *map[string]*ApiMethod {
	return &api.methods
}

func (api *Api) GetMethod(name string) *ApiMethod {
	return api.methods[name]
}

func (api *Api) GetJson() *string {
	return api.apiJson
}

//Check whether the Api is valid
func (api *Api) verify() error {
	return nil
}
