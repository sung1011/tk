package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.Hidden = true
}

var dbCmd = &cobra.Command{
	Use:   "db [玩家rid]",
	Short: "查询玩家db数据",
	Long:  `查询玩家db数据`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("需要输入rid作为参数")
		}
		return nil
	},
	Example: fmt.Sprint("tk db 1_40001\ntk db 1_40001 40002"),
	Run: func(cmd *cobra.Command, args []string) {
		rid := getRid(args[0])
		var sec string
		if len(args) == 2 {
			sec = args[1]
		} else {
			sec = getSec(rid)
		}
		sdk := getSdk(sec)
		cluster := getClusterBySdk(sdk)
		url := fmt.Sprintf("%s/index.php?&mod=web&method=Forward\\User.get&rid=%s&sec=%s&sdk_source=%s", viper.GetStringMapString("host_global")[cluster], rid, sec, sdk)
		resp, err := http.Get(url)
		if err != nil {
			log.Erro(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBufferString("")
		_, err = io.Copy(buf, resp.Body)
		log.Info(buf)
	},
}

func getRid(rid string) string {
	return rid
}
func getSec(rid string) string {
	sec := (strings.SplitAfter(rid, "_"))[1]
	if viper.IsSet("mergeSectionGroup") {
		mergeSecs := viper.GetStringMapStringSlice("mergeSectionGroup")
		for mainSec, secs := range mergeSecs {
			for _, oneSec := range secs {
				if oneSec == sec {
					return mainSec
				}
			}
		}
	}
	return sec
}

func getSdk(section string) string {
	sec, err := strconv.Atoi(section)
	if err != nil {
		log.Erro(err)
	}
	for platform, secRange := range viper.GetStringMapStringSlice("platformGroups") {
		range1, _ := strconv.Atoi(secRange[0])
		range2, _ := strconv.Atoi(secRange[1])
		if sec >= range1 && sec <= range2 {
			return platform
		}
	}
	log.Erro("bag sec " + string(sec))
	return ""
}

func getClusterBySdk(sdk string) string {
	switch sdk {
	case "ios", "ahard", "android", "adw", "amix":
		return "online"
	case "qq":
		return "qq"
	case "twonline":
		return "twonline"
	default:
		log.Erro("bag sdk " + sdk)
		return ""
	}
}
