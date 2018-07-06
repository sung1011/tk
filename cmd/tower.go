package cmd

import (
	"fmt"
	"github.com/sung1011/tk/common"
	"os/exec"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// towerCmd represents the tower command
var towerCmd = &cobra.Command{
	Use:   "tower",
	Short: "我的tower",
	Long:  `tower`,
	Example: fmt.Sprintf("%s\n%s\n%s\n",
		"tk tower me",
		"tk tower 十八寨",
		"tk tower 十八寨 回雁峰 黑龙潭",
	),
	Run: func(cmd *cobra.Command, args []string) {
		route(args)
	},
}

func init() {
	rootCmd.AddCommand(towerCmd)
	towerCmd.Hidden = true
}

func route(args []string) {
	acceptArgs := []string{"me"}
	if len(args) == 0 {
		log.Erro(fmt.Sprintf("缺少参数: %v", acceptArgs))
	}
	switch args[0] {
	case acceptArgs[0]:
		cmd := exec.Command("open", fmt.Sprintf("https://tower.im/members/%s/?me=1", viper.GetStringMapString("tower")["member"]))
		common.RunCmd(cmd)
	default:
		cmd := exec.Command("open", fmt.Sprintf("https://tower.im/teams/%s/search/?keyword=%s", viper.GetStringMapString("tower")["item"], args[0]))
		common.RunCmd(cmd)
	}
}
