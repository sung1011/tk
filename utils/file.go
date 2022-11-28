package utils

import (
	"bufio"
	"io"
	"os"
	"strings"

	log "github.com/sung1011/tk-log"
)

// Exists 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ReadLine 按行对字符串进行处理
func ReadLine(s string, handler func(string)) error {
	isfile, err := PathExists(s)
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

// ReadStdin 给与提示 + 接收stdin + 执行注册的匿名函数
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

// CopyFile copy文件
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
