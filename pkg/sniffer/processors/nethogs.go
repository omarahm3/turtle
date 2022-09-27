package processors

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

type nethogsProcessor struct {
}

func (p *nethogsProcessor) Analyze(s string) *Log {
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

func (p *nethogsProcessor) CanRunCommand() bool {
	return helpers.CommandExists("nethogs")
}

func (p *nethogsProcessor) ShouldProcess(message string) bool {
	if !strings.Contains(message, process_type) {
		return false
	}

	fields := strings.Fields(message)

	if len(fields) > 0 && fields[0] == unknown_type {
		return false
	}

	return true
}

func (p *nethogsProcessor) GetCommand() *exec.Cmd {
	return exec.Command("nethogs", "-t", "-v", "3")
}
