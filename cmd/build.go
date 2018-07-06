package cmd

import (
	"os"
	"os/exec"
	"github.com/sung1011/tk/common"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "编译",
	Long:  `编译`,
	Run: func(cmd *cobra.Command, args []string) {
		gopath := os.Getenv("GOPATH")
		c := exec.Command("go", "build", "-o", gopath+"/bin/tk", gopath+"/src/github.com/sung1011/tk/main.go")
		common.RunCmd(c)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
