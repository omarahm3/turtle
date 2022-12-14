package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const sniffer = "nethogs"

var ErrVnstatDoesNotExist = fmt.Errorf("'%s' does not exist on this machine or not exposed on current user $PATH", sniffer)

func CommandExists(c string) bool {
	_, err := exec.LookPath(c)
	return err == nil
}

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 16)
	if err != nil {
		panic(err)
	}

	return f
}

func Hashit(s string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func NormalizeString(s string) string {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	s = strings.ToLower(s)
	return strings.TrimSpace(s)
}

func ByteToMebibyte(b float64) float64 {
	return b / 1.049e+6
}

func IsRoot() bool {
	return os.Getuid() == 0
}
