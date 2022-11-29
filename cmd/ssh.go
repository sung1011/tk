package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/sung1011/tk/service"
	"github.com/sung1011/tk/utils"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sshCmd)
}

var sshCmd = &cobra.Command{
	Use:   "ssh [远程机器] [cmd]",
	Short: "ssh连接远程机器或直接执行远程命令",
	Long:  `ssh link`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("args num need 2")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sshHandler(args)
	},
}

func sshHandler(args []string) {
	var cmd string
	if len(args) == 0 {
		log.Erro("need args")
	}
	if len(args) == 1 { // no cmd
		args = append(args, "")
	} else if len(args) == 2 { // quick cmd
		switch args[1] {
		case "tail error":
			cmd = "tail -f /logs/error.log"
		}
	}
	key := args[0]
	conf := utils.GetConf("ssh", key).(map[string]interface{})
	fmt.Println("", conf["host"])
	sshConf := service.NewSSHConf(
		key,
		conf["host"].(string),
		conf["username"].(string),
		conf["password"].(string),
	)
	if err := service.RunTerminal(cmd, sshConf, os.Stdout, os.Stderr); err != nil {
		log.Erro(err)
	}
}
