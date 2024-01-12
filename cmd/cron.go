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
	Long:  `cron parse`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Erro("cron spec is required")
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
