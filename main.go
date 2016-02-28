package main

import (
	"fmt"
	"github.com/stjerncraft/controlpanelcore/config"
	"github.com/stjerncraft/controlpanelcore/core"
	"github.com/stjerncraft/controlpanelcore/services/minecraft"
	"os"
	"github.com/stjerncraft/controlpanelcore/web"
	"log"
)

func main() {
	//TODO: Check if another instance of Core is already running for this directory(File conflicts)
	//Restart self with current directory as start path, and check the command args of all running core instances for duplicates.


	//Read config
	config := config.ConfigFile("core.config", 0644)
	err := config.Read()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Found no existing config file, creating new.")
			config.Write()
		} else {
			fmt.Println("Failed to read config file: ", err)
		}
	}

	//TODO: Store the jsonApi for later comparing with API requests

	//Start Core API
	coreInst, err := core.NewCore(config)
	if err != nil {
		log.Fatal("Error while initializing core: ", err)
		return
	}

	//TODO: Initialize Users, Groups and permissions from Config

	//Core Service: Minecraft Server providers
	for _, localProvider := range *config.GetLocalMinecraftInstances() {
		localMinecraftProvider := minecraftprovider.NewMinecraftProvider(&localProvider)
		coreInst.AddServiceProvider("LocalMinecraftProvider", localMinecraftProvider, nil)
	}

	fmt.Println("Starting web server")

	//Initialize Web server
	webServer, err := web.NewWebServer(coreInst)
	if err != nil {
		log.Fatal("Error while creating web server: ", err)
		return
	}
	webServer.StartServer()
	//TODO: listen for fatal errors from Web Server and handle them(Channel)(I.e if it fails to start at all)

	//TODO: Connect to Service Cores listed in config/Local server
	//TODO: Initialize remote Service Providers
	//TODO: Request API's from Service Providers

	fmt.Println("Done")
}

