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

	"github.com/spf13/cobra"
)

const myIPAddress = "192.168.30.78"
const ipv4HeaderSize = 20

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use: "ping [ip-addr]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}

		localAddr, err := net.ResolveIPAddr("ip4", myIPAddress)
		checkError(err)
		ipAddr := args[0]
		remoteAddr, err := net.ResolveIPAddr("ip4", ipAddr)
		checkError(err)
		conn, err := net.DialIP("ip4:icmp", localAddr, remoteAddr)
		checkError(err)
		var msg [512]byte
		msg[0] = 8
		msg[1] = 0
		msg[2] = 0
		msg[3] = 0
		msg[4] = 0
		msg[5] = 13
		msg[6] = 0
		msg[7] = 37
		len := 8

		check := checkSum(msg[0:len])
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 255)

		_, err = conn.Write(msg[0:len])
		checkError(err)
		fmt.Print("Message sent:    ")
		for n := 0; n < 8; n++ {
			fmt.Print(" ", msg[n])
		}
		fmt.Println()
		// receive a reply
		size, err2 := conn.Read(msg[0:])
		checkError(err2)

		fmt.Print("Message received:")
		for n := ipv4HeaderSize; n < size; n++ {
			fmt.Print(" ", msg[n])
		}
		fmt.Println()
	},
}

func checkSum(msg []byte) uint16 {
	sum := 0
	// assume even for now
	for n := 0; n < len(msg); n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}
func init() {
	RootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
