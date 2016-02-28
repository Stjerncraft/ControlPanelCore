package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type configFile struct {
	filePath       string      //Full path to file
	filePermission os.FileMode //Mode/Permission, ex: 0644

	LocalMinecraftInstances []ConfigLocalMinecraftInstance
	ServiceCores            []ConfigServiceCore
	UserGroups              []ConfigUserGroup
	Users                   []ConfigUser
}

func ConfigFile(filePath string, filePermission os.FileMode) *configFile {
	config := new(configFile)

	config.filePath = filePath
	config.filePermission = filePermission

	config.LocalMinecraftInstances = make([]ConfigLocalMinecraftInstance, 0)
	config.ServiceCores = make([]ConfigServiceCore, 0)
	config.UserGroups = make([]ConfigUserGroup, 0)
	config.Users = make([]ConfigUser, 0)

	//config.Users = append(config.Users, ConfigUser{Name: "Wildex999", Permissions: make([]string, 0)})

	return config
}

func (config *configFile) Read() error {
	data, err := ioutil.ReadFile(config.filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	fmt.Printf("Config Read: %+v\n", config)

	return nil
}

func (config *configFile) Write() error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	//Backup existing config if it exists
	_, err = os.Stat(config.filePath)
	if err == nil {
		//Remove old backup
		backupFile := config.filePath + ".backup"
		_, err = os.Stat(backupFile)
		if err == nil {
			err = os.Remove(backupFile)
			if err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		//Backup
		err = os.Rename(config.filePath, backupFile)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(config.filePath, data, config.filePermission)

	return err
}

func (config *configFile) GetLocalMinecraftInstances() *[]ConfigLocalMinecraftInstance {
	return &config.LocalMinecraftInstances
}

func (config *configFile) SetLocalMinecraftInstances(instances *[]ConfigLocalMinecraftInstance) {
	config.LocalMinecraftInstances = *instances
}

func (config *configFile) GetServiceCores() *[]ConfigServiceCore {
	return &config.ServiceCores
}

func (config *configFile) SetServiceCores(serviceCores *[]ConfigServiceCore) {
	config.ServiceCores = *serviceCores
}

func (config *configFile) GetUserGroups() *[]ConfigUserGroup {
	return &config.UserGroups
}

func (config *configFile) SetUserGroups(userGroups *[]ConfigUserGroup) {
	config.UserGroups = *userGroups
}

func (config *configFile) GetUsers() *[]ConfigUser {
	return &config.Users
}

func (config *configFile) SetUsers(users *[]ConfigUser) {
	config.Users = *users
}
