package backends

import "os/exec"

type BackendType int

const (
	TYPE_NETHOGS BackendType = iota
	TYPE_BANDWHICH
)

type Log struct {
	App       string
	TotalPath string
	Sent      float64
	Received  float64
}

type Backend interface {
	ShouldProcess(message string) bool
	GetCommand() *exec.Cmd
	CanRunCommand() bool
	Analyze(s string) *Log
}

func New(t BackendType) Backend {
	switch t {
	case TYPE_NETHOGS:
		return newNethogsBackend()
	case TYPE_BANDWHICH:
		return newBandwhichBackend()
	default:
		return nil
	}
}

func newNethogsBackend() *nethogsBackend {
	return &nethogsBackend{}
}

func newBandwhichBackend() *bandwhichBackend {
	return &bandwhichBackend{}
}
