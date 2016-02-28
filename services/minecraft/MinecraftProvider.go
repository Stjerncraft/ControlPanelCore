package minecraftprovider

import (
	"bufio"
	"fmt"
	"github.com/stjerncraft/controlpanelcore/api"
	"github.com/stjerncraft/controlpanelcore/config"
	"github.com/stjerncraft/controlpanelcore/core"
	"github.com/stjerncraft/controlpanelcore/process"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/**
The Minecraft Service Provider will provide a list of local Minecraft servers.
*/

type MinecraftProvider struct {
	serviceProvider *core.ServiceProvider

	id         string
	location   string
	serverJar  string
	serverArgs []string
	javaPath   string
	javaArgs   []string
	autoStart  bool

	pid        int
	running    bool
	properties *MinecraftServerProperties

	processControl chan int32

	//TODO: consoleStreamIn
	consoleOutReader  io.Reader
	channelConsoleOut chan *string
	consoleErrReader  io.Reader
	channelConsoleErr chan *string
}

type MinecraftServerProperties struct {
	properties map[string]string
}

const (
	ConsoleTypeOut = "Out"
	ConsoleTypeErr = "Err"
)

type ConsoleLine struct {
	Line        *string
	ConsoleType string
}

var linesToBuffer = 64
var processManager = process.NewProcessManager()

func NewMinecraftProvider(config *config.ConfigLocalMinecraftInstance) *MinecraftProvider {
	newMinecraftProvider := &MinecraftProvider{
		id:         config.Id,
		location:   config.Location,
		serverJar:  config.ServerJar,
		serverArgs: config.ServerArgs,
		javaPath:   config.JavaPath,
		javaArgs:   config.JavaArgs,
		autoStart:  config.AutoStart,
	}

	newMinecraftProvider.processControl = make(chan int32)

	return newMinecraftProvider
}

func (self *MinecraftProvider) OnAdd(serviceProvider *core.ServiceProvider) {
	self.serviceProvider = serviceProvider

	//Add server to core, marked as status Unknown while we try to parse it

	//Lookup java processes and see if it's running
	processList, err := processManager.ListJavaProcesses()
	if err != nil {
		//Print error for Service
		fmt.Println("Error while getting process list: ", err)
		return
	}

	self.running = false
	for _, process := range processList {
		serverPath := ""
		commandLine := process.GetCommandArgs()
		if len(commandLine) > 0 {
			//Parse path from Command Line
			jarPrefix := "-jar "
			jarSuffix := ".jar"
			startPos := strings.Index(commandLine, jarPrefix)
			if startPos != -1 {
				startPos += len(jarPrefix)
				endPos := strings.Index(commandLine[startPos:], jarSuffix)
				if endPos != -1 {
					endPos += len(jarSuffix)
					serverPath = commandLine[startPos : startPos+endPos]
				}
			}
		}
		//TODO: Also check if jarFile is the correct one and give warning if not
		//serverRoot, jarFile := ...
		serverRoot, _ := filepath.Split(filepath.Clean(serverPath))
		absPath, err := filepath.Abs(serverRoot)
		if err != nil {
			fmt.Println("Error while determining Java jar path: ", err)
			return
		}
		if absPath == self.location {
			self.running = true
		}
	}

	//Read server properties
	serverProp, err := self.readServerProperties()
	if serverProp != nil {
		self.properties = serverProp
	}

	//If not running, should we start it
	if !self.running && self.autoStart {
		self.StartServer()
	}

	//Start Goroutine to watch the Server
	go self.process()

	fmt.Println("Added MinecraftProvider for: ", self.location)
}

func (self *MinecraftProvider) OnRemove() {
	close(self.processControl) //Stop the Process thread
	//TODO: Remove the core registered Server
}

func (self *MinecraftProvider) OnAddApi(api *api.Api) {

}

func (self *MinecraftProvider) OnRemoveApi(api *api.Api) {

}

func (self *MinecraftProvider) OnSessionStart() {

}

func (self *MinecraftProvider) OnSessionEnd() {

}

func (self *MinecraftProvider) OnEventRegister() error {
	return nil
}

func (self *MinecraftProvider) OnEventUnregister() {

}

func (self *MinecraftProvider) OnMethodCall(method *api.ApiMethod, args *[]api.ApiValue) api.ApiValue {
	return nil
}

func (self *MinecraftProvider) OnResourceRequest(name string, args []string) {

}

func (self *MinecraftProvider) process() {
	bufferSize := 100
	lineBuffer := make([]ConsoleLine, bufferSize)
ProcessLoop:
	for {
		select {
		case _, cont := <-self.processControl:
			if !cont {
				break ProcessLoop
			}
		default:

		}

		//Batch lines every half second
		lineCount := 0
	ConsoleLoop:
		for {
			select {
			case line := <-self.channelConsoleOut:
				if line != nil {
					lineBuffer[lineCount] = ConsoleLine{Line: line, ConsoleType: ConsoleTypeOut}
					fmt.Println("Out:", *line)
				}
			case line := <-self.channelConsoleErr:
				if line != nil {
					lineBuffer[lineCount] = ConsoleLine{Line: line, ConsoleType: ConsoleTypeErr}
					fmt.Println("Err:", *line)
				}
			default:
				break ConsoleLoop
			}
			lineCount++
			if lineCount >= bufferSize {
				break ConsoleLoop
			}
		}

		//TODO: Send event with gathered lines

		time.Sleep(100 * time.Millisecond)
	}
}

func (self *MinecraftProvider) readServerProperties() (*MinecraftServerProperties, error) {
	file, err := os.Open(filepath.Join(self.location, "server.properties"))
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error while opening server.properties: ", err)
		return nil, err
	}

	newProps := new(MinecraftServerProperties)
	newProps.properties = make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		args := strings.SplitN(line, "=", 2)
		if len(args) < 2 {
			fmt.Println("Unable to parse property: ", line)
			continue
		}
		newProps.properties[args[0]] = args[1]
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading server.properties: ", err)
		return nil, err
	}

	return newProps, nil
}

func (self *MinecraftProvider) StartServer() {
	//Setup environment
	fmt.Println("StartServer")
	ex := "java"
	if len(self.javaPath) > 0 {
		ex = filepath.Join(self.javaPath, "bin", "java")
	}
	jarPath := filepath.Join(self.location, self.serverJar)
	args := self.javaArgs
	args = append(args, "-jar", jarPath)
	args = append(args, self.serverArgs...)

	cmd := exec.Command(ex, args...)
	cmd.Dir = self.location

	//Hook into stream
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	self.consoleOutReader = stdOut
	self.channelConsoleOut = make(chan *string, linesToBuffer)
	go watchStreamLines(stdOut, self.channelConsoleOut)

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	self.consoleErrReader = stdErr
	self.channelConsoleErr = make(chan *string, linesToBuffer)
	go watchStreamLines(stdErr, self.channelConsoleErr)

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	//Start server
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	stdIn.Write([]byte("stop\n"))
	go self.process()
	cmd.Wait()
}

func watchStreamLines(streamReader io.Reader, ch chan<- *string) {
	defer close(ch)
	reader := bufio.NewReader(streamReader)
	multiLine := false
	var lineBuilder []byte
	for {
		line, isPrefix, err := reader.ReadLine()
		if !isPrefix {
			if !multiLine {
				str := string(line)
				ch <- &str
			} else {
				multiLine = false
				str := string(append(lineBuilder, line...))
				ch <- &str
			}
		} else {
			if !multiLine {
				lineBuilder = line
			} else {
				multiLine = true
				lineBuilder = append(lineBuilder, line...)
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err)
			}
			break
		}
	}
}
