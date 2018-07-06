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
	"flag"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sung1011/tk-log"
)

// scanportCmd represents the scanport command
var scanportCmd = &cobra.Command{
	Use:   "scanport",
	Short: "scanport",
	Long:  `扫描端口`,
	Run: func(cmd *cobra.Command, args []string) {
		handler()
	},
}

func handler() {
	//TODO host支持域名
	//支持多host
	host := viper.GetString("scanportHost")
	ports := viper.GetString("scanportPorts")
	port, err := strconv.Atoi(ports)
	ip := net.ParseIP(host)
	// 用于协程任务控制
	wg := sync.WaitGroup{}
	if err != nil {
		matched, _ := regexp.Match(`^\d+~\d+$`, []byte(ports))
		if !matched {
			log.Erro("bad ports", ports)
		} else {
			portSecs := strings.Split(ports, "~")
			startPort, err1 := strconv.Atoi(portSecs[0])
			endPort, err2 := strconv.Atoi(portSecs[1])
			if err1 != nil || err2 != nil || startPort < 1 || endPort < 2 || endPort <= startPort || viper.GetInt("scanportGoroutineNum") < 1 {
				flag.Usage()
			} else {
				wg.Add(endPort - startPort + 1)
				// 用于控制协程数
				parallelChan := make(chan int, viper.GetInt("scanportGoroutineNum"))
				for i := startPort; i <= endPort; i++ {
					parallelChan <- 1
					go checkPort(ip, i, &wg, &parallelChan)
				}
				wg.Wait()
			}
		}
	} else {
		handleOne(&wg, ip, port)
	}
}

func handleOne(wg *sync.WaitGroup, ip net.IP, port int) {
	wg.Add(1)
	parallelChan := make(chan int)
	go func() {
		parallelChan <- 1
	}()
	go checkPort(ip, port, wg, &parallelChan)
	wg.Wait()
}

func printOpeningPort(port int) {
	fmt.Println("port " + strconv.Itoa(port) + " is opening")
}

func checkPort(ip net.IP, port int, wg *sync.WaitGroup, parallelChan *chan int) {
	defer wg.Done()
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if err == nil {
		printOpeningPort(port)
		conn.Close()
	}
	<-*parallelChan
}

func init() {
	rootCmd.AddCommand(scanportCmd)

	scanportCmd.Flags().StringP("host", "H", "127.0.0.1", "指定域名")
	viper.BindPFlag("scanportHost", scanportCmd.Flags().Lookup("host"))
	scanportCmd.Flags().StringP("ports", "p", "80", "指定端口 如:80, 20000~210000")
	viper.BindPFlag("scanportPorts", scanportCmd.Flags().Lookup("ports"))
	scanportCmd.Flags().IntP("goroutinenum", "n", 10, "协程数")
	viper.BindPFlag("scanportGoroutineNum", scanportCmd.Flags().Lookup("goroutinenum"))
}
