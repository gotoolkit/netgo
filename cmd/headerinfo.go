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
	"net"
	"os"

	"io/ioutil"

	"github.com/spf13/cobra"
)

// headerinfoCmd represents the headerinfo command
var headerinfoCmd = &cobra.Command{
	Use: "headerinfo [service]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}
		service := args[0]
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		checkError(err)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		checkError(err)
		_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
		checkError(err)
		result, err := ioutil.ReadAll(conn)
		checkError(err)

		fmt.Println(string(result))
	},
}

func init() {
	RootCmd.AddCommand(headerinfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// headerinfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// headerinfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
