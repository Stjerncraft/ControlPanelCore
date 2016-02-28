package process

type Thread struct {
}

type ProcessManager interface {
	ListProcesses() ([]Process, error)
	ListJavaProcesses() ([]Process, error)
	//ListThreads(process Process) ([]Thread, error)
}
