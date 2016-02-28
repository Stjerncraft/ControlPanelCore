package api

import "errors"

type ApiObject struct {
	api  *Api
	name string

	inherit   string            //Name of parent object
	ancestors map[string]string //Map of all ancestors

	hasInherited bool
	fields       map[string]ApiType

	//GetField(name string) ApiType
	//GetFields() map[string]ApiType
}

func NewApiObject(name string, api *Api) *ApiObject {
	obj := new(ApiObject)
	obj.name = name
	obj.api = api

	return obj
}

func (obj *ApiObject) SetFields(fields *map[string]string) error {
	obj.inherit = (*fields)["inherit"]

	obj.fields = make(map[string]ApiType)
	for fieldName, value := range *fields {
		newType, err := ApiTypeFromString(obj.api, value)
		if err != nil {
			return err
		}
		obj.fields[fieldName] = newType
	}
	delete(obj.fields, "inherit")

	return nil
}

func (obj *ApiObject) InheritFields() error {
	if obj.hasInherited {
		return nil
	}
	obj.hasInherited = true

	if len(obj.inherit) == 0 {
		return nil
	}

	parentObj := obj.api.GetObject(obj.inherit)
	if parentObj == nil {
		return errors.New("Could not find parent Object: " + obj.inherit + " while inheriting fields for " + obj.name)
	}

	parentObj.InheritFields() //Propagate the inherit call up the chain

	for f, t := range parentObj.fields {
		_, exists := obj.fields[f] //Parent can't overwrite child defined type
		if !exists {
			obj.fields[f] = t
		}
	}

	return nil
}
