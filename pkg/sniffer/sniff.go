package sniffer

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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

func Sniff(sl chan SniffLog) {
	message := make(chan string)
	cmd := exec.Command("nethogs", "-v", mb_collect, "-t")

	var wg sync.WaitGroup
	wg.Add(1)

	go runCmd(cmd, &wg, message)
	go listen(message, sl)

	wg.Wait()
}

func listen(message chan string, sl chan SniffLog) {
	for {
		m := <-message
		if !strings.Contains(m, process_type) {
			continue
		}

		fields := strings.Fields(m)

		if len(fields) > 0 && fields[0] == unknown_type {
			continue
		}

		// debug(fields)
		l := analyze(fields)
		sl <- l
	}
}

func analyze(output []string) SniffLog {
	app := output[0]
	totalPath := strings.Join(output[:len(output)-2], "")
	received := output[len(output)-1]
	sent := output[len(output)-2]

	return SniffLog{
		App:           app,
		AppHash:       hashit(app),
		TotalPath:     totalPath,
		TotalPathHash: hashit(totalPath),
		Sent:          parseFloat(sent),
		Received:      parseFloat(received),
	}
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 16)
	if err != nil {
		panic(err)
	}

	return f
}

func hashit(s string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func debug(output []string) {
	fmt.Println("---------------------------")
	for i := 0; i < len(output); i++ {
		fmt.Printf("[%d] = %q\n", i, output[i])
	}
	fmt.Println("---------------------------")
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
		// fmt.Print(s.Text())
		message <- s.Text()
	}
}
