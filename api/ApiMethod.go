package api

type ApiMethod struct {
	Name string

	Args   []ApiType
	Return ApiType
}

func NewApiMethod(name string, args []string, ret string, api *Api) (*ApiMethod, error) {
	newMethod := new(ApiMethod)

	newMethod.Name = name
	newMethod.Args = make([]ApiType, len(args))

	//Method Arguments
	for index, typ := range args {
		argType, err := ApiTypeFromString(api, typ)
		if err != nil {
			return nil, err
		}

		newMethod.Args[index] = argType
	}

	//Method return type
	retType, err := ApiTypeFromString(api, ret)
	if err != nil {
		return nil, err
	}
	newMethod.Return = retType

	return newMethod, nil
}
