package cmd

import (
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sung1011/tk-log"

	"github.com/jinzhu/now"

	"github.com/spf13/cobra"
)

var timeArg string

const (
	timeFLagNow = iota
	timeFlagST
	timeFlagTS
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "转化时间",
	Long:  `转化时间`,
	Run: func(cmd *cobra.Command, args []string) {
		switch chkArgs(args) {
		case timeFlagST:
			strtotime()
		case timeFlagTS:
			timetostr()
		}
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
}

func chkArgs(args []string) int {
	if len(args) == 0 || args[0] == "now" {
		_now()
		os.Exit(0)
	}
	timeArg = strings.Join(args, " ")
	if strings.Count(timeArg, "")-1 == 10 && strings.Index(timeArg, ":") == -1 && strings.Index(timeArg, "-") == -1 && strings.Index(timeArg, "/") == -1 {
		return timeFlagTS
	}
	return timeFlagST
}

func strtotime() {
	loc, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		log.Erro(err)
	}
	t, err := now.ParseInLocation(loc, timeArg)
	if err != nil {
		log.Erro(t, err)
	}
	log.Succ(t.Unix())
	log.Succ(t.Format("2006-01-02 15:04:05"))
}

func timetostr() {
	i, err := strconv.ParseInt(timeArg, 10, 64)
	if err != nil {
		log.Erro(i, err)
	}
	t := time.Unix(i, 0)
	log.Succ(t.Unix())
	log.Succ(t.Format("2006-01-02 15:04:05"))
}

func _now() {
	t := time.Now()
	log.Succ(t.Unix())
	log.Succ(t.Format("2006-01-02 15:04:05"))
}
