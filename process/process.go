package process

type Process interface {
	GetId() uint32
	GetParentId() uint32
	GetName() string
	GetExecutablePath() string
	GetCommandArgs() string
}
