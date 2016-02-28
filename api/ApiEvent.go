package api

type ApiEvent struct {
	args []ApiType
}

func NewApiEvent(args *[]string, api *Api) (*ApiEvent, error) {
	newApiEvent := new(ApiEvent)
	newApiEvent.args = make([]ApiType, len(*args))
	for index, val := range *args {
		valType, err := ApiTypeFromString(api, val)
		if err != nil {
			return nil, err
		}
		newApiEvent.args[index] = valType
	}

	return newApiEvent, nil
}
