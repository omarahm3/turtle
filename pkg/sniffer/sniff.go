package sniffer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/omarahm3/turtle/pkg/helpers"
	"github.com/omarahm3/turtle/pkg/sniffer/processors"
)

var (
	processor processors.Processor
)

const (
	unknown_type = "unknown"
	process_type = "/"
	mb_collect   = "3"
)

type AppLog struct {
	App      string  `json:"app"`
	AppHash  string  `json:"appHash"`
	Sent     float64 `json:"sent"`
	Received float64 `json:"received"`
}

func (a AppLog) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func ToAppLog(s string) AppLog {
	var a AppLog
	err := json.Unmarshal([]byte(s), &a)
	if err != nil {
		panic(err)
	}

	return a
}

type SniffLog struct {
	App           string  `json:"app"`
	AppHash       string  `json:"appHash"`
	TotalPath     string  `json:"totalPath"`
	TotalPathHash string  `json:"hash"`
	Sent          float64 `json:"sent"`
	Received      float64 `json:"received"`
}

func (s SniffLog) NewAppLog() AppLog {
	return AppLog{
		App:      s.App,
		AppHash:  s.AppHash,
		Sent:     s.Sent,
		Received: s.Received,
	}
}

func (s SniffLog) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func ToSniffLog(s string) SniffLog {
	var l SniffLog
	err := json.Unmarshal([]byte(s), &l)
	if err != nil {
		panic(err)
	}

	return l
}

func Sniff(sl chan SniffLog, p processors.Processor) {
	processor = p
	message := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)

	go runCmd(processor.GetCommand(), &wg, message)
	go listen(message, sl)

	wg.Wait()
}

func listen(message chan string, sl chan SniffLog) {
	for {
		m := <-message
		if !processor.ShouldProcess(m) {
			continue
		}

		l := processor.Analyze(m)
		sl <- SniffLog{
			App:           l.App,
			AppHash:       helpers.Hashit(l.App),
			TotalPath:     l.TotalPath,
			TotalPathHash: helpers.Hashit(l.TotalPath),
			Sent:          l.Sent,
			Received:      l.Received,
		}
	}
}

func runCmd(cmd *exec.Cmd, wg *sync.WaitGroup, message chan string) {
	defer wg.Done()
	sout, _ := cmd.StdoutPipe()
	go readStd(sout, message)

	serr, _ := cmd.StderrPipe()
	go readStd(serr, message)

	err := cmd.Start()
	if err != nil {
		fmt.Printf("run:: error occurred running command: %q", err.Error())
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("run:: error occurred running command: %q", err.Error())
		return
	}
}

func readStd(r io.ReadCloser, message chan string) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		message <- s.Text()
	}
}
