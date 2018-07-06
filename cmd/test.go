package cmd

import (
	"fmt"
	"github.com/sung1011/tk/common"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "just for test",
	Long:  `just for test`,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("ls", "-al", ".")
		common.RunCmd(c)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Hidden = true
}

func loading() {
	fn := func(arg interface{}) {
		time.Sleep(time.Millisecond * 10000)
		fmt.Println(arg)
	}
	common.Loading(fn, "a")
}
