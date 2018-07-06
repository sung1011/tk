package service

import (
	"io"
	"os"

	log "github.com/sung1011/tk-log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

//SSHConf ssh config
type SSHConf struct {
	tag  string
	host string
	un   string
	pw   string
}

// NewSSHConf pass
func NewSSHConf(tag, host, un, pw string) *SSHConf {
	return &SSHConf{
		tag, host, un, pw,
	}
}

//RunTerminal ssh登录, 执行命令
func RunTerminal(cmd string, conf *SSHConf, stdout, stderr io.Writer) error {
	//connection
	conn, err := ssh.Dial("tcp", conf.host, &ssh.ClientConfig{
		User:            conf.un,
		Auth:            []ssh.AuthMethod{ssh.Password(conf.pw)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	log.NewPreface().SetMulti(map[string]string{
		"username": conf.un,
		"tag":      conf.tag,
	}).Show()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	if cmd != "" {
		session.Run(cmd)
	} else {
		session.Shell()
		session.Wait()
	}
	return nil
}

//RunCmd 仅执行ssh命令
func RunCmd(conf SSHConf, cmd string, stdout, stderr io.Writer) {
	//connection
	conn, err := ssh.Dial("tcp", conf.host, &ssh.ClientConfig{
		User:            conf.un,
		Auth:            []ssh.AuthMethod{ssh.Password(conf.pw)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Erro(err)
	}
	defer conn.Close()
	//session
	session, err := conn.NewSession()
	if err != nil {
		log.Erro(err)
	}
	defer session.Close()

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin
	// run
	if err := session.Run(cmd); err != nil {
		log.Erro(err)
	}
}
