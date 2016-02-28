package json

import "encoding/json"

type JsonApi struct {
	Name string

	Objects map[string]JsonApiObject
	Events  map[string]JsonEventObject
	Methods map[string]JsonApiMethod
}

//An object is just a list of Field name to Type
type JsonApiObject map[string]string

type JsonEventObject []string

type JsonApiMethod struct {
	Args   []string
	Return string
}

func ParseJsonApi(jsonString string) (*JsonApi, error) {
	jsonApi := new(JsonApi)
	err := json.Unmarshal([]byte(jsonString), jsonApi)

	return jsonApi, err
}

func GetJsonApi(jsonApi *JsonApi) (string, error) {
	json, err := json.Marshal(jsonApi)

	return string(json), err
}

//Check whether the two API's are the same
func (api *JsonApi) Equals(otherApi *JsonApi) bool {
	return false
}
