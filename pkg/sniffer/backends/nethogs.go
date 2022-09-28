package backends

import (
	"os/exec"
	"strings"

	"github.com/omarahm3/turtle/pkg/helpers"
)

const (
	unknown_type = "unknown"
	process_type = "/"
	mb_collect   = "3"
)

type nethogsBackend struct {
}

func (n *nethogsBackend) Analyze(s string) *Log {
	output := strings.Fields(s)
	app := helpers.NormalizeString(output[0])
	totalPath := helpers.NormalizeString(strings.Join(output[:len(output)-2], ""))
	received := helpers.NormalizeString(output[len(output)-1])
	sent := helpers.NormalizeString(output[len(output)-2])

	return &Log{
		App:       app,
		TotalPath: totalPath,
		Received:  helpers.ParseFloat(received),
		Sent:      helpers.ParseFloat(sent),
	}
}

func (n *nethogsBackend) CanRunCommand() bool {
	return helpers.CommandExists("nethogs")
}

func (n *nethogsBackend) ShouldProcess(message string) bool {
	if !strings.Contains(message, process_type) {
		return false
	}

	fields := strings.Fields(message)

	if len(fields) > 0 && fields[0] == unknown_type {
		return false
	}

	return true
}

func (n *nethogsBackend) GetCommand() *exec.Cmd {
	return exec.Command("nethogs", "-t", "-v", "3")
}
