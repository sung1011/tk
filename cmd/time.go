package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-module/carbon/v2"
	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "转化时间",
	Long:  `转化时间`,
	Run: func(cmd *cobra.Command, args []string) {
		var cb carbon.Carbon
		var value string
		var tz string = carbon.PRC
		if len(args) == 0 {
			value = "now"
		} else {
			value = strings.Join(args, " ")
		}
		v, err := strconv.Atoi(value)
		if err == nil {
			// 时间戳
			// v, _ := strconv.Atoi(value)
			// fmt.Println("", value, v)
			if len(value) == 10 {
				cb = cb.CreateFromTimestamp(int64(v), tz)
			}
			if len(value) == 13 {
				cb = cb.CreateFromTimestampMilli(int64(v), tz)
			}
		} else {
			// 字符串
			cb = cb.Parse(value, tz)
		}
		log.Succ(fmt.Sprintf("时区: %s", tz))
		log.Succ(fmt.Sprintf("字符: %s", cb.ToDateTimeString()))
		log.Succ(fmt.Sprintf("时间戳: %d", cb.Carbon2Time().Unix()))
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
