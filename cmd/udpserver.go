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

// udpserverCmd represents the udpserver command
var udpserverCmd = &cobra.Command{
	Use: "udpserver",
	Run: func(cmd *cobra.Command, args []string) {
		service := ":1200"
		udpAddr, err := net.ResolveUDPAddr("udp", service)
		checkError(err)
		conn, err := net.ListenUDP("udp", udpAddr)
		checkError(err)
		for {
			handleUDPClient(conn)
			fmt.Println("Listen & Serve")
		}
	},
}

func handleUDPClient(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}

func init() {
	RootCmd.AddCommand(udpserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// udpserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// udpserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
