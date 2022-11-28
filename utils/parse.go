package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/sung1011/tk-log"
	"github.com/tidwall/gjson"
)

// Value 从一段带有冒号(:)的字符串中 解析出值
func Value(s, k string) string {
	strs := strings.Split(s, ":")
	if k != "" && !strings.Contains(strs[0], k) {
		return ""
	}
	if len(strs) >= 2 {
		return strings.Trim(strs[1], " ")
	}
	return ""
}

// Loading async show load progress
// TODO: go fn
// ISSUE: *File.Sync()
func Loading(fn func(interface{}), arg interface{}) {
	tik := time.NewTicker(time.Millisecond * 1000)
	go func() {
		i := 0
		for _ = range tik.C {
			i++
			fmt.Printf("\rLoading%s", strings.Repeat(".", i%4))
			os.Stdout.Sync()
		}
	}()
	fn(arg)
	tik.Stop()
}

// RunCmd 执行命令
func RunCmd(cmd *exec.Cmd) {
	var args string
	if len(cmd.Args) > 0 {
		for _, arg := range cmd.Args[1:] {
			args += " " + arg
		}
	}
	log.ShowPreface("cmd", cmd.Path+" "+args)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	color.Set(color.FgCyan)
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Erro(fmt.Sprintf("cmd.Start() failed with '%s'\n", err))
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, stdoutIn)
	if err != nil {
		log.Erro(err)
	}
	str := buf.String()
	if !gjson.Valid(str) {
		_, errStdout = io.Copy(stdout, bytes.NewReader(buf.Bytes()))
	} else {
		var outBuf bytes.Buffer
		err := json.Indent(&outBuf, buf.Bytes(), "", "\t")
		if err != nil {
			log.Erro("marshalIndent error", err)
		}
		_, errStdout = io.Copy(stdout, bytes.NewReader(outBuf.Bytes()))
	}
	_, errStderr = io.Copy(stderr, stderrIn)
	err = cmd.Wait()
	if err != nil {
		log.Erro(fmt.Sprintf("cmd.Run() failed with %s\n", err))
	}
	if errStdout != nil || errStderr != nil {
		log.Erro("failed to capture stdout or stderr\n")
	}
}
