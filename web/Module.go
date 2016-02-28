package web

const (
	AccessTypePrivate = 0
	AccessTypePublic = 1
)
type AccessType int

type Module struct {
	name string
	access AccessType
	description string
	//Permissions
}

func NewModule(name string, access AccessType) *Module {
	newModule := new(Module)
	newModule.name = name
	newModule.access = access

	return newModule
}

func (m *Module) getName() string {
	return m.name
}

func (m *Module) getAccess() AccessType {
	return m.access
}
