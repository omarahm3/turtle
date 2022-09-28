package backends

import (
	"os/exec"
	"strings"

	"github.com/omarahm3/turtle/pkg/helpers"
)

type bandwhichBackend struct {
}

func (b *bandwhichBackend) Analyze(s string) *Log {
	output := strings.Fields(s)
	output = output[2:]
	app := helpers.NormalizeString(output[0])
	output = output[3:]
	metrics := strings.Split(output[0], "/")
	sent := helpers.ParseFloat(helpers.NormalizeString(metrics[0]))
	received := helpers.ParseFloat(helpers.NormalizeString(metrics[1]))

	return &Log{
		App:       app,
		TotalPath: app,
		Received:  helpers.ByteToMebibyte(received),
		Sent:      helpers.ByteToMebibyte(sent),
	}
}

func (b *bandwhichBackend) CanRunCommand() bool {
	return helpers.CommandExists(b.GetCommand().Args[0])
}

func (b *bandwhichBackend) ShouldProcess(message string) bool {
	if !strings.Contains(message, "process:") {
		return false
	}
	return true
}

func (b *bandwhichBackend) GetCommand() *exec.Cmd {
	return exec.Command("bandwhich", "-t", "-p", "--raw")
}
