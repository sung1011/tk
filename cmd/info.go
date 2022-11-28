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
	"go/build"
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

// RuntimeInfo hold info
type RuntimeInfo struct {
	GOVersion string
	GOOS      string
	GOARCH    string
	NumCPU    int
	GOPATH    string
	GOROOT    string
	Compiler  string
	PATH      string `split:":"`
	SHELL     string
	CurPath   string
}

func info() {
	ri := RuntimeInfo{
		GOVersion: getGOVersion(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
		NumCPU:    runtime.NumCPU(),
		GOPATH:    build.Default.GOPATH,
		GOROOT:    runtime.GOROOT(),
		Compiler:  runtime.Compiler,
		PATH:      os.Getenv("PATH"),
		SHELL:     os.Getenv("SHELL"),
		CurPath:   getCurPath(),
	}
	t := reflect.TypeOf(ri)
	v := reflect.ValueOf(ri)
	for i := 0; i < t.NumField(); i++ {
		splitSeq, exists := t.Field(i).Tag.Lookup("split")
		var val interface{}
		if exists {
			val = strings.Split(v.Field(i).Interface().(string), splitSeq)
		} else {
			val = v.Field(i).Interface()
		}
		log.Info(fmt.Sprintf("%-9v : %v", t.Field(i).Name, val))
	}
}

func getGOVersion() string {
	v, err := exec.Command("go", "version").Output()
	if err != nil {
		log.Erro(err)
	}
	return strings.Split(string(v), " ")[2]
}

func getCurPath() string {
	d, err := os.Getwd()
	if err != nil {
		log.Erro(err)
	}
	return d
}
