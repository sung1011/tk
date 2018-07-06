// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	log "github.com/sung1011/tk-log"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "计算机基础信息",
	Long:  `info`,
	Run: func(cmd *cobra.Command, args []string) {
		info()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

//RuntimeInfo hold info
type RuntimeInfo struct {
	GOVersion string
	GOOS      string
	GOARCH    string
	NumCPU    int
	GOPATH    string
	GOROOT    string
	Compiler  string
	PATH      string
	SHELL     string
	HOME      string
}

func info() {
	ri := RuntimeInfo{
		getGOVersion(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		os.Getenv("GOPATH"),
		runtime.GOROOT(),
		runtime.Compiler,
		os.Getenv("PATH"),
		os.Getenv("SHELL"),
		getHOME(),
	}

	t := reflect.TypeOf(ri)
	v := reflect.ValueOf(ri)
	for i := 0; i < t.NumField(); i++ {
		log.Info(fmt.Sprintf("%-9v : %v", t.Field(i).Name, v.Field(i).Interface()))
	}

}

func getGOVersion() string {
	v, err := exec.Command("go", "version").Output()
	if err != nil {
		log.Erro(err)
	}
	return strings.Split(string(v), " ")[2]
}

func getHOME() string {
	d, err := os.Getwd()
	if err != nil {
		log.Erro(err)
	}
	return d
}
