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
	"math"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	fixInterest = 1
	fixCapital  = 2
)

// loanCmd represents the loan command
var loanCmd = &cobra.Command{
	Use:   "loan",
	Short: "loan贷款计算器",
	Long:  `loan`,
	Run: func(cmd *cobra.Command, args []string) {
		Handler()
	},
}

// Handler 处理
func Handler() {
	switch viper.GetInt("loanType") {
	case fixInterest:
		l := &Loan{
			money:        viper.GetFloat64("loanMoney"),
			ratePerMonth: viper.GetFloat64("loanRate") / 12,
			month:        viper.GetFloat64("loanYear") * 12,
			adv:          viper.GetFloat64("loanAdvance"),
			t:            fixInterest,
		}
		l.monthlySupply()
	case fixCapital:

	}
}

// Loan 基础结构
type Loan struct {
	money        float64 `name:"贷款总额"`
	ratePerMonth float64 `name:"月利息"`
	month        float64 `name:"还款月数"`
	adv          float64 `name:"每月提前还款额"`
	t            int     `name:"还款方式(1.等息 2.等金)"`
}

// monthSupply 处理月供
// 等额本息
// [贷款本金×月利率×（1+月利率）^还款月数]÷[（1+月利率）^还款月数－1]
// 等额本金
// [贷款本金×月利率×（1+月利率）^还款月数]÷[（1+月利率）^还款月数－1]
func (l *Loan) monthlySupply() {
	r := l.money * l.ratePerMonth * math.Pow((1+l.ratePerMonth), l.month) / (math.Pow((1+l.ratePerMonth), l.month) - 1)
	fmt.Println(l)
	fmt.Println("-----")
	fmt.Println(r)
}

func init() {
	rootCmd.AddCommand(loanCmd)
	loanCmd.Hidden = true

	loanCmd.Flags().IntP("ty", "t", 1, "贷款方式")
	loanCmd.Flags().Float64P("money", "v", 0.0, "贷款额度")
	loanCmd.Flags().Float64P("year", "y", 25, "年限")
	loanCmd.Flags().Float64P("rate", "r", 0, "利率")
	loanCmd.Flags().Float64P("adv", "a", 0, "提前还款额")

	viper.BindPFlag("loanType", loanCmd.Flags().Lookup("ty"))
	viper.BindPFlag("loanMoney", loanCmd.Flags().Lookup("money"))
	viper.BindPFlag("loanYear", loanCmd.Flags().Lookup("year"))
	viper.BindPFlag("loanRate", loanCmd.Flags().Lookup("rate"))
	viper.BindPFlag("loanAdvance", loanCmd.Flags().Lookup("adv"))
}
