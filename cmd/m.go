package cmd

import (
	"fmt"
	"github.com/sung1011/tk/common"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	log "github.com/sung1011/tk-log"
)

// mCmd represents the m command
var mCmd = &cobra.Command{
	Use:   "m",
	Short: "master2接口调试",
	Long:  `master2接口调试`,
	Example: fmt.Sprintf("%s\n%s\n%s\n",
		"tk m 291_d1 User.get",
		"tk m 291_d1 PVP\\Arena.get",
		"tk m 291_d1 Debug.test foo=bar",
	),
	Run: func(cmd *cobra.Command, args []string) {
		mHandler(args)
	},
}

func init() {
	rootCmd.AddCommand(mCmd)
	mCmd.Hidden = true
}

func mHandler(args []string) {
	if len(args) < 2 {
		log.Erro(fmt.Sprintf("缺少必要参数: m {rid} {method}"))
	}
	rid, sec, method, params := mHandlArgs(args)
	url := fmt.Sprintf("http://127.0.0.1:8080/connector_backend/index.php?&mod=web&sec=%s&method=%s&version=dev&rid=%s%s", sec, method, rid, params)
	cmd := exec.Command("curl", "-s", url)
	common.RunCmd(cmd)
}

func mHandlArgs(args []string) (string, string, string, string) {
	rid := args[0]
	sec := (strings.SplitAfter(rid, "_"))[1]
	method := args[1]
	params := ""
	if len(args) > 2 {
		for _, a := range args[2:] {
			params += "&" + a
		}
	} else {
		params = ""
	}
	return rid, sec, method, params
}
