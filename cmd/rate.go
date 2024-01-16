package cmd

import (
	"fmt"
	"strconv"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
)

var rateCmd = &cobra.Command{
	Use:   "rate 352 10000 2",
	Short: "rate parse",
	Long:  `param: cur, total, costTime(min)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			log.Erro("args num need 3")
		}
		cur, err := strconv.Atoi(args[0])
		if err != nil {
			log.Erro(err)
		}
		total, err := strconv.Atoi(args[1])
		if err != nil {
			log.Erro(err)
		}
		costMin, err := strconv.Atoi(args[2])
		if err != nil {
			log.Erro(err)
		}
		if cur > total {
			log.Erro("cur > total")
		}
		speed := float32(cur) / float32(costMin)
		needTime := float32(total-cur) * float32(costMin) / float32(cur)

		log.Succ(fmt.Sprintf("进度: %.2f%%", float32(cur)*100/float32(total)))
		log.Succ(fmt.Sprintf("速度: %.2f/min | %.2f/hour", speed, speed*60))
		log.Succ(fmt.Sprintf("需要时间: %.2fmin | %.2fhour", needTime, needTime/60))
	},
}

func init() {
	rootCmd.AddCommand(rateCmd)
}
