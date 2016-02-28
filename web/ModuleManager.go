package web

import "errors"

type ModuleManager struct {
	modules map[string]*Module
}

func NewModuleManager() *ModuleManager {
	newManager := new(ModuleManager)

	newManager.modules = make(map[string]*Module)

	return newManager
}

//Create a new module and add it. Returns the new module. If there already exists a module with the given name,
//this will fail with an error.
func (manager *ModuleManager) NewModule(name string, access AccessType) (*Module, error) {
	if manager.modules[name] != nil {
		return nil, errors.New("Module with the given name already exists: " + name)
	}

	newModule := NewModule(name, access)
	manager.modules[name] = newModule

	return newModule, nil
}

//Add the given module. If a module with the same name already exists, it will be overwritten, and the original returned.
func (manager *ModuleManager) AddModule(module *Module) *Module {
	oldModule := manager.modules[module.name]
	manager.modules[module.name] = module

	return oldModule
}

//Returns nil if no module with the given name exists
func (manager *ModuleManager) GetModule(name string) *Module {
	return manager.modules[name]
}


