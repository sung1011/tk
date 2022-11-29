package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-module/carbon/v2"
	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
)

var (
	tz string
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "time converter",
	Long:  `time converter`,
	Run: func(cmd *cobra.Command, args []string) {
		var cb carbon.Carbon
		var value string
		if len(args) == 0 {
			value = "now"
		} else {
			value = strings.Join(args, " ")
		}
		v, err := strconv.Atoi(value)
		if err == nil {
			if len(value) == 10 { // 时间戳
				cb = cb.CreateFromTimestamp(int64(v), tz)
			}
			if len(value) == 13 { // 时间戳 毫秒
				cb = cb.CreateFromTimestampMilli(int64(v), tz)
			}
		} else {
			// 字符串
			cb = cb.Parse(value, tz)
		}
		log.Succ(fmt.Sprintf(" %s", tz))
		log.Succ(fmt.Sprintf(" %s", cb.ToDateTimeString()))
		log.Succ(fmt.Sprintf(" %d", cb.Carbon2Time().Unix()))
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
	timeCmd.PersistentFlags().StringVar(&tz, "tz", carbon.PRC, "set timezone")
}
