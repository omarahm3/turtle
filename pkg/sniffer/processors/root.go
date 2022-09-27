package processors

import "os/exec"

type ProcessorType int

const (
	TYPE_NETHOGS ProcessorType = iota
	TYPE_BANDWHICH
)

type Log struct {
	App       string
	TotalPath string
	Sent      float64
	Received  float64
}

type Processor interface {
	ShouldProcess(message string) bool
	GetCommand() *exec.Cmd
	CanRunCommand() bool
	Analyze(s string) *Log
}

func New(t ProcessorType) Processor {
	switch t {
	case TYPE_NETHOGS:
		return newNethogsProcessor()
	case TYPE_BANDWHICH:
		return newBandwhichProcessor()
	default:
		return nil
	}
}

func newNethogsProcessor() *nethogsProcessor {
	return &nethogsProcessor{}
}

func newBandwhichProcessor() *bandwhichProcessor {
	return &bandwhichProcessor{}
}
