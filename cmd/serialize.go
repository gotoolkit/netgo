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
	"encoding/gob"
	"os"

	"github.com/spf13/cobra"
)

type Person struct {
	Name  Name
	Email []Email
}
type Name struct {
	Family   string
	Personal string
}
type Email struct {
	Kind    string
	Address string
}

// serializeCmd represents the serialize command
var serializeCmd = &cobra.Command{
	Use: "serialize",
	Run: func(cmd *cobra.Command, args []string) {
		person := Person{

			Name: Name{Family: "Newmarch", Personal: "Jan"},
			Email: []Email{
				Email{Kind: "home", Address: "jan@newmarch.name"},
				Email{Kind: "work", Address: "j.newmarch@boxhill.edu.au"},
			},
		}
		saveGob("person.gob", person)
	},
}

func saveGob(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
	outFile.Close()
}

func init() {
	RootCmd.AddCommand(serializeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serializeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serializeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
