package processors

import (
	"os/exec"
	"strings"

	"github.com/omarahm3/turtle/pkg/helpers"
)

type bandwhichProcessor struct {
}

func (p *bandwhichProcessor) Analyze(s string) *Log {
	output := strings.Fields(s)
	output = output[2:]
	app := helpers.NormalizeString(output[0])
	output = output[3:]
	metrics := strings.Split(output[0], "/")
	sent := helpers.NormalizeString(metrics[0])
	received := helpers.NormalizeString(metrics[1])

	return &Log{
		App:      app,
		Received: helpers.ParseFloat(received),
		Sent:     helpers.ParseFloat(sent),
	}
}

func (p *bandwhichProcessor) CanRunCommand() bool {
	return helpers.CommandExists(p.GetCommand().Args[0])
}

func (p *bandwhichProcessor) ShouldProcess(message string) bool {
	if !strings.Contains(message, "process:") {
		return false
	}
	return true
}

func (p *bandwhichProcessor) GetCommand() *exec.Cmd {
	return exec.Command("bandwhich", "-t", "-p", "--raw")
}
