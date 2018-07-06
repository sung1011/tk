package cmd

import (
	"errors"
	"os"

	"github.com/sung1011/tk/service"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(sshCmd)
}

var sshCmd = &cobra.Command{
	Use:   "ssh [远程机器tag]",
	Short: "ssh连接远程机器或直接执行远程命令",
	Long:  `ssh link`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("参数只要2个")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		runTerStd(args)
	},
}

func runTerStd(args []string) {
	cmd, conf := preRun(args)
	if err := service.RunTerminal(cmd, conf, os.Stdout, os.Stderr); err != nil {
		log.Erro(err)
	}
}

func preRun(args []string) (string, *service.SSHConf) {
	var cmd string
	if len(args) == 0 {
		args = append(args, "board")
	}
	// alias := map[string]string{"d": "dev", "dev": "dev", "md": "mat-dev", "mat-dev": "mat-dev", "b": "board", "board": "board"}
	// viper.Set("sshmap", alias)
	// if v, ok := alias[args[0]]; !ok {
	// 	log.Erro("参数错误:", v, "可选参数:", alias)
	// }
	// args[0] = alias[args[0]]
	if len(args) == 1 {
		// 无cmd 表达ssh登录
		args = append(args, "")
	} else if len(args) == 2 {
		//快捷操作
		switch args[1] {
		case "tail -f error":
			cmd = ""
		}
	} else {
		log.Erro("参数个数最多2个")
	}
	sshconf := service.NewSSHConf(
		args[0],
		viper.GetStringMapString("host_ssh")[args[0]],
		viper.GetStringMapString("username")[args[0]],
		viper.GetStringMapString("password")[args[0]],
	)
	return cmd, sshconf
}
