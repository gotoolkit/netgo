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
	"time"

	"github.com/spf13/cobra"
)

// daytimeserverCmd represents the daytimeserver command
var daytimeserverCmd = &cobra.Command{
	Use: "daytimeserver",
	Run: func(cmd *cobra.Command, args []string) {
		service := ":1200"
		tcpAddr, err := net.ResolveTCPAddr("tcp", service)
		checkError(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		checkError(err)
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			if multiThread {
				go handleTCPClient(conn)
			} else {
				daytime := time.Now().String()
				conn.Write([]byte(daytime))
				conn.Close()
			}
			fmt.Println("Listen & Serve")
		}
	},
}

func handleTCPClient(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:]))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

var multiThread bool

func init() {
	RootCmd.AddCommand(daytimeserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daytimeserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daytimeserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	daytimeserverCmd.Flags().BoolVarP(&multiThread, "multi-thread-mode", "m", false, "enable Multi-Thread Server")

}
