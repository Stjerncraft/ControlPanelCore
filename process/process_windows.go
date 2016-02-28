package process

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type process struct {
	id             uint32
	parentId       uint32
	name           string
	executablePath string

	commandArgsStored bool
	commandArgs       string
}

func (p *process) GetId() uint32 {
	return p.id
}

func (p *process) GetParentId() uint32 {
	return p.parentId
}

func (p *process) GetName() string {
	return p.name
}

func (p *process) GetExecutablePath() string {
	return p.executablePath
}

//Warning: This takes up to a second on windows, use sparingly.
//Warning2: Trying to run in threads causes problems
func (p *process) GetCommandArgs() string {
	if p.commandArgsStored {
		return p.commandArgs
	}

	cmd := exec.Command("wmic", "process", "where", "ProcessId="+strconv.FormatUint(uint64(p.id), 10), "get", "Commandline")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to get command line arguments for process: " + strconv.FormatUint(uint64(p.id), 10))
		fmt.Println(err)
		return ""
	} else {
		outString := out.String()

		prefix := "CommandLine"
		indexPrefix := strings.Index(outString, prefix)
		if indexPrefix < 0 || len(prefix)+1 >= len(outString) {
			fmt.Println("Failed to get Command Line Arguments, got: " + outString)
			return ""
		}
		p.commandArgs = strings.TrimSpace(outString[len(prefix):])
	}

	p.commandArgsStored = true
	return p.commandArgs
}
