package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sung1011/tk/common"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	log "github.com/sung1011/tk-log"

	"github.com/jinzhu/gorm"
	//sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tidwall/gjson"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "获取线上log (会被缓存到本地)",
	Long: fmt.Sprint(
		"表名固定rows\n",
		"schema固定，别名只可以是a,b,c,d\n",
		"驼峰字段需要转下划线 errorMsg -> error_msg\n",
	),
	Example: fmt.Sprint(
		"tk log -q \"select optype,count(*) as a from rows group by optype\"\n",
		"tk log -q \"select id,time from rows limit 3\"\n",
		"tk log -q \"select error_message from rows limit 1\"\n",
	),
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Hidden = true

	logCmd.Flags().StringP("cluster", "c", "online", "指定集群")
	viper.BindPFlag("clusterFlag", logCmd.Flags().Lookup("cluster"))

	now := time.Now()
	date := fmt.Sprintf("%d%.2d%.2d", now.Year(), now.Month(), now.Day())
	logCmd.Flags().StringP("date", "d", date, "指定日期")
	viper.BindPFlag("date", logCmd.Flags().Lookup("date"))

	logCmd.Flags().BoolP("force", "f", false, "强制拉取log")
	viper.BindPFlag("logforce", logCmd.Flags().Lookup("force"))

	logCmd.Flags().StringP("mod", "m", "backend", "指定模块 (backend, global)")
	viper.BindPFlag("logmod", logCmd.Flags().Lookup("mod"))

	logCmd.Flags().StringP("type", "t", "backend", "指定细节类型 (backend, task, sdk_backend)")
	viper.BindPFlag("logtype", logCmd.Flags().Lookup("type"))

	// logCmd.Flags().StringP("query", "q", "", "查询 (使用的gjson库 https://github.com/tidwall/gjson)")
	logCmd.Flags().StringP("query", "q", "", "查询 (sql)")
	viper.BindPFlag("query", logCmd.Flags().Lookup("query"))
}

func run() {
	cluster := viper.GetString("clusterFlag")
	date := viper.GetString("date")
	file := "all_" + viper.GetString("logtype") + ".txt"
	tmppath := "/tmp"
	fullpath := tmppath + "/" + cluster + "_" + viper.GetString("logmod") + "_" + date + "_" + file

	ex, err := common.Exists(fullpath)
	if err != nil {
		log.Erro(err)
	}
	if ex == false || viper.GetBool("logforce") == true {
		host := viper.GetStringMapString("host_admin")
		writefile(host[cluster]+"/"+viper.GetString("logmod")+"/"+date+"/"+file, fullpath)
	}
	if viper.GetString("query") != "" {
		// import file to struct
		rows := toRows(fullpath)
		// import data to sqlite
		toDB(rows)
		// sqlite query
		query(viper.GetString("query"))
	}
}

func writefile(url, path string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Erro(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Erro(err)
	}
	if err = ioutil.WriteFile(path, b, 0644); err != nil {
		log.Erro(err)
	}
}

//Row 数据行
type Row struct {
	UID string `gorm:"primary_key;AUTO_INCREMENT"`

	ID           string `json:"id,omitempty"`
	Optype       string `json:"optype,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Type         string `json:"type,omitempty"`
	File         string `json:"file,omitempty"`
	Line         string `json:"line,omitempty"`
	Level        string `json:"level,omitempty"`
	Rid          string `json:"rid,omitempty"`
	Sec          string `json:"sec,omitempty"`
	Time         string `json:"time,omitempty"`
	Version      string `json:"version,omitempty"`
	Trace        string `json:"trace,omitempty"`
	Request      string `json:"request,omitempty"`
	SystemStatus string `json:"systemStatus,omitempty"`
	UserData     string `json:"userData,omitempty"`
	ErrorData    string `json:"errorData,omitempty"`

	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
	C string `json:"c,omitempty"`
	D string `json:"d,omitempty"`
}

func toRows(path string) []Row {
	lineNumOneRow := 17
	var rows []Row
	var i int
	row := new(Row)
	var lineHandler = func(l string) {
		if len(l) == 0 {
			return
		}
		if i >= lineNumOneRow {
			rows = append(rows, *row)
			row = new(Row)
			i = 0
		}
		i++
		sl := strings.Split(l, ":")
		k := ""
		if sl[0] == "uid" {
			k = "UID"
		} else if sl[0] == "id" {
			k = "ID"
		} else if strings.TrimSpace(sl[0]) == "<html>" {
			log.Erro("network error, try again by -f")
		} else {
			k = strings.Title(sl[0])
		}
		v := strings.Join(sl[1:], "")
		rvs := reflect.ValueOf(row).Elem()
		rvs.FieldByName(k).SetString(v)
	}
	common.ReadLine(path, lineHandler)
	return rows
}

func toJSON(rows []Row) string {
	jb, err := json.Marshal(rows)
	if err != nil {
		log.Erro(err)
	}
	return string(jb)
}

func gJSON(j string) {
	rs := gjson.Get(j, viper.GetString("query"))
	log.Succ(rs)
}

func toDB(rows []Row) {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		log.Erro(err, db)
	}
	defer db.Close()
	db.AutoMigrate(&Row{})
	if viper.GetBool("logforce") {
		db.DropTable(&Row{})
	}
	if !db.HasTable(&Row{}) {
		db.CreateTable(&Row{})
		for _, row := range rows {
			if db.NewRecord(row) {
				db.Create(&row)
			}
		}
	}
}

func query(sql string) {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		log.Erro(err, db)
	}
	defer db.Close()
	db = db.Table("rows")
	var rs []Row
	db.Raw(sql).Scan(&rs)
	for _, v := range rs {
		b, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			log.Erro(err)
		}
		log.Trac(string(b))
	}
}
