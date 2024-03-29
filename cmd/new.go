package cmd

import (
	"os/exec"

	"github.com/sung1011/tk/utils"

	"github.com/spf13/cobra"
	log "github.com/sung1011/tk-log"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "新建一个命令",
	Long:  `新建一个命令`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Erro("缺少参数: 新建命令名")
		}
		for _, s := range args {
			c := exec.Command("cobra-cli", "add", "-t", "github.com/sung1011/tk", s)
			utils.RunCmd(c)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

}
