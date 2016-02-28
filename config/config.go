package config

type ConfigLocalMinecraftInstance struct {
	Id         string   //Used by services and modules when working against a specific server
	Name       string   //The name given for this server
	Location   string   //The folder which the Minecraft instance lives in
	ServerJar  string   //The server JAR file
	ServerArgs []string //Any Arguments to add to the Minecraft server
	JavaPath   string   //Sets JAVA_HOME to this before starting, leave empty to use existing environment variable
	JavaArgs   []string //Arguments for the Java VM
	string              //Any arguments to add to the Java VM
	AutoStart  bool     //Whether or not the server should to try to start this Instance
}

type ConfigServiceCore struct {
	Address string
	Port    uint16
}

type ConfigUserGroup struct {
	Name string

	ParentId uint32

	Permissions []string
}

type CoreService struct {
	serviceName string
	enabled     bool
}

type ConfigUser struct {
	Name string

	Groups      []ConfigUserGroup
	Permissions []string

	CoreServices []CoreService
}

type Config interface {
	Read() error  //Read the Config from it's source
	Write() error //Write to storage

	//List of local Minecraft servers
	GetLocalMinecraftInstances() *[]ConfigLocalMinecraftInstance
	SetLocalMinecraftInstances(instances *[]ConfigLocalMinecraftInstance)

	//List of Service Cores
	GetServiceCores() *[]ConfigServiceCore
	SetServiceCores(serviceCores *[]ConfigServiceCore)

	//User Groups
	GetUserGroups() *[]ConfigUserGroup
	SetUserGroups(userGroups *[]ConfigUserGroup)

	//Users
	GetUsers() *[]ConfigUser
	SetUsers(users *[]ConfigUser)
}
