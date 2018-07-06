package common

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sung1011/tk-log"
	"github.com/tidwall/gjson"

	"github.com/fatih/color"
)

//Value 从一段带有冒号(:)的字符串中 解析出值
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

//Exists 路径是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//ChkExists 路径是否存在处理
func ChkExists(path string) {
	ex, err := Exists(path)
	if err != nil {
		log.Erro(err)
	}
	if !ex {
		log.Erro("目录不存在" + path)
	}
}

//ReadLine 按行对字符串进行处理
func ReadLine(s string, handler func(string)) error {
	isfile, err := Exists(s)
	if err != nil {
		return err
	}
	var f io.Reader
	if isfile {
		f, err = os.Open(s)
		if err != nil {
			return err
		}
	} else {
		f = strings.NewReader(s)
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

//ReadStdin 给与提示 + 接收stdin + 执行注册的匿名函数
func ReadStdin(prompt string, handler func(string) bool) (input string) {
	log.Warn(prompt)
	buf := bufio.NewReader(os.Stdin)
	for {
		bs, _, _ := buf.ReadLine()
		input = strings.TrimSpace(string(bs))
		if handler(input) {
			break
		}
	}
	return input
}

//CopyFile copy文件
func CopyFile(source, dest string) {
	if source == "" || dest == "" {
		log.Erro("empty param")
	}
	sourceOpen, err := os.Open(source)
	if err != nil {
		log.Erro(err)
	}
	defer sourceOpen.Close()
	destOpen, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Erro(err)
	}
	defer destOpen.Close()
	_, copyErr := io.Copy(destOpen, sourceOpen)
	if copyErr != nil {
		log.Erro(err)
	}
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

//RunCmd 执行命令
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
