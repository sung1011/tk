package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "cron parse",
	Long:  `cron 秒 分 小时 日期 月份 星期几 年份`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Erro(`格式: <cron 秒 分 小时 日期 月份 星期几 年份>

　　（1）*：表示匹配该域的任意值。假如在Minutes域使用*, 即表示每分钟都会触发事件。
　　（2）?：只能用在DayofMonth和DayofWeek两个域。它也匹配域的任意值，但实际不会。因为DayofMonth和DayofWeek会相互影响。例如想在每月的20日触发调度，不管20日到底是星期几，则只能使用如下写法： 13 13 15 20 * ?, 其中最后一位只能用？，而不能使用*，如果使用*表示不管星期几都会触发，实际上并不是这样。
　　（3）-：表示范围。例如在Minutes域使用5-20，表示从5分到20分钟每分钟触发一次 
　　（4）/：表示起始时间开始触发，然后每隔固定时间触发一次。例如在Minutes域使用5/20,则意味着5分钟触发一次，而25，45等分别触发一次. 
　　（5）,：表示列出枚举值。例如：在Minutes域使用5,20，则意味着在5和20分每分钟触发一次。 

	0 15 10 * * ?     每天上午10:15触发
	0 0/30 9-17 * * ?   朝九晚五工作时间内每半小时
	0 0 10,14,16 * * ?   每天上午10点，下午2点，4点 
	0 0 12 ? * WED    表示每个星期三中午12点
			`)
		}
		s := strings.Join(args, " ")
		cron, err := ParseCron(s)
		if err != nil {
			log.Erro(err)
		}
		t1 := cron.Next(time.Now())
		log.Succ(t1.Format("2006-01-02 15:04:05"))
		t2 := cron.Next(t1)
		log.Succ(t2.Format("2006-01-02 15:04:05"))
		t3 := cron.Next(t2)
		log.Succ(t3.Format("2006-01-02 15:04:05"))
		t4 := cron.Next(t3)
		log.Succ(t4.Format("2006-01-02 15:04:05"))
		t5 := cron.Next(t4)
		log.Succ(t5.Format("2006-01-02 15:04:05"))
		t6 := cron.Next(t5)
		log.Succ(t6.Format("2006-01-02 15:04:05"))
		t7 := cron.Next(t6)
		log.Succ(t7.Format("2006-01-02 15:04:05"))
		t8 := cron.Next(t7)
		log.Succ(t8.Format("2006-01-02 15:04:05"))
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}

type Cron struct {
	schedule *cron.SpecSchedule
}

var parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func ParseCron(spec string) (*Cron, error) {
	schedule, err := parser.Parse(spec)
	if err != nil {
		return nil, err
	}
	if specSchedule, ok := schedule.(*cron.SpecSchedule); !ok {
		return nil, fmt.Errorf("bad cron %s", spec)
	} else {
		return &Cron{schedule: specSchedule}, nil
	}
}

func (c *Cron) Next(t time.Time) time.Time {
	return c.schedule.Next(t)
}
