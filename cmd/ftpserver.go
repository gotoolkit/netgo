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
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// ftpserverCmd represents the ftpserver command
var ftpserverCmd = &cobra.Command{
	Use: "ftpserver",
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
			go handleFTPClient(conn)
		}
	},
}

const (
	DIR = "DIR"
	CD  = "CD"
	PWD = "PWD"
)

func handleFTPClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			conn.Close()
			return
		}
		s := string(buf[0:n])
		if s[0:2] == CD {
			cddir(conn, s[3:])
		} else if s[0:3] == DIR {
			dirList(conn)
		} else if s[0:3] == PWD {
			pwd(conn)
		}
	}
}

func dirList(conn net.Conn) {
	defer conn.Write([]byte("\r\n"))
	dir, err := os.Open(".")
	if err != nil {
		return
	}
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}
	for _,nm := range names {
		conn.Write([]byte(nm + "\r\n"))
	}
}

func pwd(conn net.Conn) {
	s, err := os.Getwd()
	if err != nil {
		conn.Write([]byte(""))
		return
	}
	conn.Write([]byte(s))
}

func cddir(conn net.Conn, s string) {
	log.Println(s)
	err := os.Chdir(s)
	if err == nil {
		conn.Write([]byte("OK"))
	} else {
		conn.Write([]byte(err.Error()))
	}
}

func init() {
	RootCmd.AddCommand(ftpserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ftpserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ftpserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
