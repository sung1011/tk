package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sung1011/tk/common"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cluster 集群
type argsType []string

var aupMap = map[string]bool{
	"IOS_Stable":      true,
	"Android_Stable":  true,
	"Android_Default": true,
	"IOS_Default":     true,
}

func init() {
	rootCmd.AddCommand(chkfcCmd)
	chkfcCmd.Hidden = true

	chkfcCmd.Flags().StringP("cluster", "c", "online", "指定集群")
	viper.BindPFlag("chkfccluster", chkfcCmd.Flags().Lookup("cluster"))
}

var chkfcCmd = &cobra.Command{
	Use:     "chkfc [线上集群]",
	Short:   "获取线上filecache文件",
	Long:    `check filecache`,
	Example: fmt.Sprintf("%s\n%s\n%s\n%s", "tk chkfc (default online)", "tk chkfc online", "tk chkfc qq", "tk chkfc twonline"),
	Run: func(cmd *cobra.Command, args []string) {
		cluster := getCluster(args)
		host := viper.GetStringMapString("host_admin")
		resp, err := http.Get(host[cluster] + "/filecache/meta/version.yaml")
		if err != nil {
			log.Erro(err)
		}
		defer resp.Body.Close()
		buf := bytes.NewBufferString("")
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			log.Erro(err)
		}
		r := bufio.NewReader(buf)
		n := 0
		upgrade := ""
		buf1 := bytes.NewBuffer([]byte(""))
		gap := ""
		for {
			l, _, eof := r.ReadLine()
			line := strings.Trim(string(l), " ")
			up := common.Value(line, "upgrade_path")
			if up != "" || n > 1 {
				buf1 = bytes.NewBuffer([]byte(""))
				upgrade = up
				n = 0
			} else {
				n++
			}
			if n > 1 {
				log.Info(line)
				continue
			}
			if n == 0 {
				llen := len(line)
				offset := 40
				if llen < offset {
					gap = strings.Repeat(" ", offset-llen)
				}
				buf1.WriteString(line + gap)
			} else {
				buf1.WriteString("\t" + line)
			}
			v := common.Value(line, "version")
			if v != "" {
				// 根据升级序列 查询详情
				resp, err = http.Get(host[cluster] + "/versioninfo/" + v + "/info.yaml")
				if err != nil {
					log.Erro(err)
				}
				defer resp.Body.Close()
				buf2 := bytes.NewBufferString("")
				_, err = io.Copy(buf2, resp.Body)
				if err != nil {
					log.Erro(err)
				}
				rr := bufio.NewReader(buf2)
				// 将详情折叠存储
				var tmp, info1, info2, info3, info4, info5 string
				for {
					ll, _, innerEOF := rr.ReadLine()
					tmp = common.Value(string(ll), "backend_tag")
					if tmp != "" {
						info1 = tmp
					}
					tmp = common.Value(string(ll), "global_tag")
					if tmp != "" {
						info2 = tmp
					}
					tmp = common.Value(string(ll), "config_tag")
					if tmp != "" {
						info3 = tmp
					}
					tmp = common.Value(string(ll), "script_tag")
					if tmp != "" {
						info4 = tmp
					}
					tmp = common.Value(string(ll), "vms_tag")
					if tmp != "" {
						info5 = tmp
					}
					if innerEOF == io.EOF {
						break
					}
				}
				info := fmt.Sprintf("\t%s %s %s %s %s", info1, info2, info3, info4, info5)
				buf1.WriteString(info)
				// 是否重要升级序列
				isMainUp := false
				for _, mainup := range viper.GetStringSlice("upgradepath") {
					if mainup == upgrade {
						log.Succ(buf1)
						isMainUp = true
					}

				}
				if isMainUp == false {
					log.Info(buf1)
				}
			}
			if eof == io.EOF {
				break
			}
		}
	},
}

func getCluster(args []string) string {
	c := ""
	if len(args) == 0 {
		c = viper.GetString("chkfccluster")
	} else {
		c = args[0]
	}
	cs := viper.GetStringMapString("cluster")
	for k := range cs {
		if k == c {
			return c
		}
	}
	log.Erro("cluster参数错误")
	os.Exit(1)
	return ""
}
