package cmd

import (
	"fmt"
	"github.com/sung1011/tk/common"
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/sung1011/tk-log"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	hotfixConf map[string]interface{}
)

// hotfixCmd represents the hotfix command
var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "交互模式给线上打hotfix",
	Long:  `hotfix`,
	Example: fmt.Sprintf("%s\n%s\n%s\n",
		"tk hotfix",
		"tk hotfix -m config",
		"tk hotfix -m backend -t 18.7.1 -f",
	),
	Run: func(cmd *cobra.Command, args []string) {
		getConf()
		if viper.GetString("hxmod") == "" {
			common.ReadStdin("指定模块 ["+strings.Join(viper.GetStringSlice("mods"), " ")+"]:", setMod)
		}

		//
		if viper.GetString("tag") != "" {
			hotfixConf["tag"] = viper.GetString("tag")
			diffFile(hotfixConf["tag"].(string))
		} else {
			if hotfixConf["mod"].(string) == "config" {
				hotfixConf["tag"] = common.ReadStdin("输入线上tag(推荐chkfc查看):", pass)
				common.ReadStdin("指定文件(多个以空格分隔):", setDiffFiles)
			} else {
				hotfixConf["tag"] = common.ReadStdin("输入线上tag(推荐chkfc查看):", diffFile)
			}
		}

		//
		if viper.GetBool("hotfixforce") {
			diffCode("Y")
		} else {
			common.ReadStdin("是否显示与线上("+hotfixConf["tag"].(string)+")差异(y/n):", diffCode)
		}

		cpFiles()

		showParam()

		common.ReadStdin("是否上传(y/n):", isUpload)

		openTianti()

		tartf()

		log.Succ("Done~")
	},
}

func pass(v string) bool {
	return true
}
func getConf() {
	m := make(map[string]interface{})

	m["proj"] = viper.GetString("proj")
	m["mod"] = viper.GetString("hxmod")
	m["cluster"] = viper.GetString("hotfixcluster")
	m["branch"] = viper.GetString("branch")
	m["projPath"] = viper.GetStringMapString("pathway")["project"] + "/" + m["proj"].(string) + "/" + "version/" + m["mod"].(string) + "/" + m["branch"].(string)
	m["hotfixAppPath"] = viper.GetStringMapString("pathway")["hotfix_app"] //git@gitlab.playcrab-inc.com:wangtao02/upload_hotfix.git
	m["hotfixTmpPath"] = m["hotfixAppPath"].(string) + "/" + "tmp"
	m["diffFiles"] = []string{}

	common.ChkExists(m["projPath"].(string))
	common.ChkExists(m["hotfixAppPath"].(string))

	hotfixConf = m
}

func setMod(v string) bool {
	if len(v) == 0 {
		v = viper.GetString("hxmod")
	}
	mods := viper.GetStringSlice("mods")
	for _, mod := range mods {
		if mod == v {
			hotfixConf["mod"] = v
			hotfixConf["projPath"] = viper.GetStringMapString("pathway")["project"] + "/" + hotfixConf["proj"].(string) + "/" + "version/" + hotfixConf["mod"].(string) + "/" + hotfixConf["branch"].(string)
			return true
		}
	}
	log.Erro("mod参数错误")
	return false
}

func setDiffFiles(v string) bool {
	if len(v) == 0 {
		return false
	}
	hotfixConf["diffFiles"] = strings.Fields(v)
	return true
}

func diffFile(v string) bool {
	//git pull
	// cmd := exec.Command("git", "pull")
	// cmd.Dir = hotfixConf["projPath"].(string)
	// b, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Erro(err, string(b))
	// }
	// log.Warn(cmd.Args)
	// log.Info(string(b))

	//git diff file
	cmd := exec.Command("git", "diff", "--name-only", v)
	cmd.Dir = hotfixConf["projPath"].(string)
	// common.RunCmd(cmd)
	b, _ := cmd.CombinedOutput()
	log.ShowPreface("cmd", cmd.Args)
	log.Info(string(b))

	err := common.ReadLine(string(b), func(s string) {
		dfs := hotfixConf["diffFiles"].([]string)
		if len(s) != 0 {
			dfs = append(dfs, s)
		}
		hotfixConf["diffFiles"] = dfs
	})
	if err != nil {
		log.Warn(err)
	}

	return true
}

func diffCode(yn string) bool {
	if yn == "Y" || yn == "y" {
		//git diff version
		cmd := exec.Command("git", "diff", hotfixConf["tag"].(string))
		cmd.Dir = hotfixConf["projPath"].(string)
		common.RunCmd(cmd)
		// b, _ := cmd.CombinedOutput()
		// log.ShowPreface("running", cmd.Args)
		// log.Info(string(b))
	}
	return true
}

func cpFiles() {
	log.ShowPreface("running", "文件复制到临时目录 等待打压缩包")
	os.RemoveAll(hotfixConf["hotfixTmpPath"].(string))
	for _, f := range hotfixConf["diffFiles"].([]string) {
		sourDir := viper.GetStringMapString("pathway")["project"] + "/" + "master" + "/version/" + hotfixConf["mod"].(string) + "/dev"
		destDir := hotfixConf["hotfixTmpPath"].(string) + "/playcrab/" + hotfixConf["proj"].(string) + "/version/" + hotfixConf["mod"].(string) + "/" + hotfixConf["tag"].(string)
		if err := os.MkdirAll(destDir+"/"+path.Dir(f), 0755); err != nil {
			log.Erro(err)
		}
		common.CopyFile(sourDir+"/"+f, destDir+"/"+f)
		log.Info(destDir + "/" + f)
	}
}

func isUpload(yn string) bool {
	if yn != "Y" && yn != "y" {
		return true
	}
	cmd := exec.Command("bash", "fixbug.sh", hotfixConf["proj"].(string), hotfixConf["cluster"].(string), hotfixConf["hotfixTmpPath"].(string))
	cmd.Dir = hotfixConf["hotfixAppPath"].(string)
	common.RunCmd(cmd)
	// log.Warn(cmd.Args)
	// b, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Erro(err, string(b))
	// }
	// log.Info(string(b))

	return true
}

func showParam() {

	log.NewPreface().SetMulti(
		map[string]string{
			"项目名":  hotfixConf["proj"].(string),
			"模块名":  hotfixConf["mod"].(string),
			"集群名":  hotfixConf["cluster"].(string),
			"线上分支": hotfixConf["tag"].(string),
		},
	).Show()
}

func openTianti() {
	cmd := exec.Command("open", viper.GetStringMapString("bookmark")["tianti"])
	common.RunCmd(cmd)
}

func tartf() {
	cmd := exec.Command("tar", "-tf", "/data/home/user00/package/hotfix.tar.gz")
	common.RunCmd(cmd)
}

func init() {
	rootCmd.AddCommand(hotfixCmd)
	hotfixCmd.Hidden = true

	hotfixCmd.Flags().StringP("proj", "p", "master", "指定项目")
	viper.BindPFlag("proj", hotfixCmd.Flags().Lookup("proj"))

	hotfixCmd.Flags().StringP("mod", "m", "backend", "指定模块")
	viper.BindPFlag("hxmod", hotfixCmd.Flags().Lookup("mod"))

	hotfixCmd.Flags().StringP("branch", "b", "dev", "指定分支")
	viper.BindPFlag("branch", hotfixCmd.Flags().Lookup("branch"))

	hotfixCmd.Flags().StringP("tag", "t", "", "指定线上tag")
	viper.BindPFlag("tag", hotfixCmd.Flags().Lookup("tag"))

	hotfixCmd.Flags().StringP("cluster", "c", "online", "指定集群")
	viper.BindPFlag("hotfixcluster", hotfixCmd.Flags().Lookup("cluster"))

	hotfixCmd.Flags().BoolP("force", "f", false, "无须询问(全yes)")
	viper.BindPFlag("hotfixforce", hotfixCmd.Flags().Lookup("force"))
}
