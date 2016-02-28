package core

const (
	MinecraftStatus_Unknown = "Unknown"
	MinecraftStatus_Offline = "Offline"
	MinecraftStatus_Online  = "Online"
)

type MinecraftServer struct {
	//Id must be an unique string which does not change
	Id string

	Name          string
	Address       string
	Status        string
	PlayersOnline int32
	PlayersMax    int32

	//Set by core
	ProviderId int32 //The id of the ServiceProvider managing this Minecraft Server
}

var CoreMinecraftServerApiJson = `{
	"name": "MinecraftServerProvider v1.0.0",

	"objects": {
		"MinecraftServer": {
			"Id": "string",
			"Name": "string",
			"Address": "string",
			"Status": "string",
			"PlayersOnline": "int",
			"PlayersMax": "int",
			"ProviderId": "int"
		},
		"ConsoleLine": {
			"type": "string",
			"line": "string"
		},
		"ConsoleOutput": {
			"lines": "list[ConsoleLine]"
		},
		"ConsoleInput": {
			"input": "string"
		}
	},
	"events": {
		"ConsoleOut": ["ConsoleOutput"],
		"ConsoleIn": ["ConsoleInput"],
		"UpdatedServer": ["MinecraftServer"]
	},
	"methods": {
		"startServer": {
			"args": [],
			"return": "byte"
		},
		"stopServer": {
			"args": [],
			"return": "byte"
		},
		"getServer": {
			"args": [],
			"return": "MinecraftServer"
		},
		"sendConsoleInput": {
			"args": ["string"],
			"return": ""
		}
	}
}`
