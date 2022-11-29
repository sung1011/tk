/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sung1011/tk-log"
	"github.com/sung1011/tk/utils"
)

// multiCmd represents the multi command
var multiCmd = &cobra.Command{
	Use:   "multi",
	Short: "multi task",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confMulti := viper.GetStringMapStringSlice("multi")
		var keys string
		for k := range confMulti {
			keys += k + " "
		}
		if len(args) == 0 {
			log.Erro("可选参数值 " + keys)
		}
		for _, arg := range args {
			multiHandler(arg)
		}
	},
}

func init() {
	rootCmd.AddCommand(multiCmd)
}

func multiHandler(arg string) {
	conf := utils.GetConf("multi", arg).([]interface{})
	var wg sync.WaitGroup
	for _, task := range conf {
		t := task.(string)
		wg.Add(1)
		go func(task string) {
			defer wg.Done()
			sp := strings.Split(task, " ")
			c := exec.Command(sp[0], sp[1:]...)
			utils.RunCmd(c)
		}(t)
	}
	wg.Wait()
}
