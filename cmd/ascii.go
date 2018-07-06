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
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var col [128]doc

// asciiCmd represents the ascii command
var asciiCmd = &cobra.Command{
	Use:   "ascii",
	Short: "ASCII码表",
	Long:  `ascii table list`,
	Run: func(cmd *cobra.Command, args []string) {
		gen()
		list()
	},
}

func gen() {
	for k := range col {
		d := doc{dec: k}
		d.toBin()
		d.toHex()
		d.toCharacter()
		col[k] = d
	}
}

type doc struct {
	dec  int    `tag:"十进制"`
	bin  string `tag:"二进制"`
	hex  string `tag:"十六进制"`
	char string `tag:"图形"`
}

func (d *doc) toBin() {
	d.bin = fmt.Sprintf("%b", d.dec)
}
func (d *doc) toHex() {
	d.hex = fmt.Sprintf("%X", d.dec)
}
func (d *doc) toCharacter() {
	d.char = fmt.Sprintf("%q", d.dec)
}

func list() {
	tmp := doc{}
	bin, _ := reflect.TypeOf(tmp).FieldByName("bin")
	dec, _ := reflect.TypeOf(tmp).FieldByName("dec")
	hex, _ := reflect.TypeOf(tmp).FieldByName("hex")
	char, _ := reflect.TypeOf(tmp).FieldByName("char")
	fmt.Printf("|%5s|%3s|%4s|%4s|\n", bin.Tag.Get("tag"), dec.Tag.Get("tag"), hex.Tag.Get("tag"), char.Tag.Get("tag"))
	fmt.Println("---------------------------------")
	allRet := ""
	indexRet := ""
	index := viper.GetString("asciiIndex")
	for _, doc := range col {
		s := fmt.Sprintf("|%8s|%6d|%8s|%6s|\n", doc.bin, doc.dec, doc.hex, doc.char)
		if strings.Contains(s, index) {
			indexRet += s
		}
		allRet += s
	}
	if indexRet != "" || index != "" {
		fmt.Println(indexRet)
	} else {
		fmt.Println(allRet)
	}
}

func init() {
	rootCmd.AddCommand(asciiCmd)

	asciiCmd.Flags().StringP("index", "i", "", "匹配")
	viper.BindPFlag("asciiIndex", asciiCmd.Flags().Lookup("index"))
}
