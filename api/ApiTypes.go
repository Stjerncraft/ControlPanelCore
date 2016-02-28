package api

import (
	"errors"
	"strings"
)

type ApiType interface {
	GetTypeAsString() string
	New() ApiValue
}

//Byte, 8 bit
type ApiTypeByte struct{}

func (t *ApiTypeByte) GetTypeAsString() string {
	return "byte"
}
func (t *ApiTypeByte) New() ApiValue {
	return new(ApiValueByte)
}

//Short, 16 bit
type ApiTypeShort struct{}

func (t *ApiTypeShort) GetTypeAsString() string {
	return "short"
}
func (t *ApiTypeShort) New() ApiValue {
	return new(ApiValueShort)
}

//Integer, 32 bit
type ApiTypeInteger struct{}

func (t *ApiTypeInteger) GetTypeAsString() string {
	return "int"
}
func (t *ApiTypeInteger) New() ApiValue {
	return new(ApiValueInteger)
}

//Long, 64 bit
type ApiTypeLong struct{}

func (t *ApiTypeLong) GetTypeAsString() string {
	return "long"
}
func (t *ApiTypeLong) New() ApiValue {
	return new(ApiValueLong)
}

//Float, 32 bit
type ApiTypeFloat struct{}

func (t *ApiTypeFloat) GetTypeAsString() string {
	return "float"
}
func (t *ApiTypeFloat) New() ApiValue {
	return new(ApiValueFloat)
}

//Double, 64 bit
type ApiTypeDouble struct{}

func (t *ApiTypeDouble) GetTypeAsString() string {
	return "double"
}
func (t *ApiTypeDouble) New() ApiValue {
	return new(ApiValueDouble)
}

//String, variable length
type ApiTypeString struct{}

func (t *ApiTypeString) GetTypeAsString() string {
	return "string"
}
func (t *ApiTypeString) New() ApiValue {
	return new(ApiValueString)
}

//Object, variable length
type ApiTypeObject struct {
	//Used to look up the ApiObject
	Name string
	Api  *Api
}

func (t *ApiTypeObject) GetTypeAsString() string {
	return t.Name
}
func (t *ApiTypeObject) New() ApiValue {
	objValue := new(ApiValueObject)
	objValue.Type = t

	return objValue
}

//List, variable length
type ApiTypeList struct {
	Type ApiType
}

func (t *ApiTypeList) GetTypeAsString() string {
	return "list[" + t.Type.GetTypeAsString() + "]"
}
func (t *ApiTypeList) New() ApiValue {
	list := new(ApiValueList)
	list.Type = t
	return list
}

func ApiTypeFromString(api *Api, str string) (ApiType, error) {
	switch str {
	case "byte":
		return new(ApiTypeByte), nil
	case "short":
		return new(ApiTypeShort), nil
	case "int":
		return new(ApiTypeInteger), nil
	case "long":
		return new(ApiTypeLong), nil
	case "float":
		return new(ApiTypeFloat), nil
	case "double":
		return new(ApiTypeDouble), nil
	case "string":
		return new(ApiTypeString), nil
	}

	if strings.HasPrefix(str, "list[") && strings.HasSuffix(str, "]") {
		prefixLen := 5
		suffixLen := 1

		if len(str) <= prefixLen+suffixLen {
			return nil, errors.New("List does not define a type: " + str)
		}

		listTypeStr := str[prefixLen : len(str)-suffixLen]
		listType, err := ApiTypeFromString(api, listTypeStr)
		if err != nil {
			return nil, err
		}

		list := new(ApiTypeList)
		list.Type = listType

		return list, nil
	}

	//If not a base type, then it must be an object type
	if api.GetObject(str) == nil {
		return nil, errors.New("Unknown type or object: " + str)
	}

	objType := new(ApiTypeObject)
	objType.Name = str
	objType.Api = api

	return objType, nil
}
