// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

	"github.com/spf13/cobra"
	"strconv"
	"net"
	"os"
)

// ipmaskCmd represents the ipmask command
var ipmaskCmd = &cobra.Command{
	Use: "ipmask",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fatal("Usage: ipmask dotted-ip-addr ones bits")
		}
		dotAddr := args[0]

		ones, _ := strconv.Atoi(args[1])
		bits, _ := strconv.Atoi(args[2])
		addr := net.ParseIP(dotAddr)

		if addr == nil {
			fmt.Println("Invalid address")
			os.Exit(1)
		}
		mask := net.CIDRMask(ones, bits)
		network := addr.Mask(mask)
		fmt.Println("Address is ", addr.String())
		fmt.Println("Mask length is ", bits)
		fmt.Println("Leading ones count is ", ones)
		fmt.Println("Mask is (HEX) ", mask.String())
		fmt.Println("Network is ", network.String())
	},
}

func init() {
	RootCmd.AddCommand(ipmaskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipmaskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipmaskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}