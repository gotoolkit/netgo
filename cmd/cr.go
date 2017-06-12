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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

type Result struct {
	SERVICOREF     string `json:"SERVICOREF"`
	MERCHREF       string `json:"MERCH_REF"`
	TXNID          string `json:"TXN_ID"`
	TXNSTATUS      string `json:"TXN_STATUS"`
	TXNRESPCODE    string `json:"TXN_RESP_CODE"`
	TXNMESSAGE     string `json:"TXN_MESSAGE"`
	TXNDATETIME    string `json:"TXN_DATE_TIME"`
	TXNAMTDEDUCTED string `json:"TXN_AMT_DEDUCTED"`
	AUTHID         string `json:"AUTH_ID"`
	POSMERCHSHOPID string `json:"POS_MERCH_SHOP_ID"`
	RISKSCORE      string `json:"RISK_SCORE"`
	Sign           string `json:"sign"`
}

// crCmd represents the cr command
var crCmd = &cobra.Command{
	Use:   "cr",
	Short: "concurrent request",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}
		url := args[0]

		if !debug {
			log.SetOutput(ioutil.Discard)
		}

		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		//start := time.Now()
		ch := make(chan string)
		go func() {
			for {
				select {
				case <-ch:
				default:
					for i := 0; i < frequency; i++ {
						go MakeRequest(url, ch)
					}
					time.Sleep(10 * time.Second)
				}

			}
		}()

		handleSignals()

	},
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

var client *http.Client

func MakeRequest(url string, ch chan string) {
	start := time.Now()
	resp, _ := client.Get(url)
	secs := time.Since(start).Seconds()

	log.Println(secs, " elapsed with response ", url)
	if resp.StatusCode == 302 {
		rUrl, _ := resp.Location()
		MakeFormRequest(rUrl.String(), ch)
	}

}
func MakeFormRequest(url string, ch chan string) {

	resp, _ := client.Get(url)
	resp.Request.ParseForm()
	form := resp.Request.Form
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	//log.Println(url)
	//log.Println(form)

	MakeJsonRequest("https://epayweb01.macaucep.gov.mo/CEPPGDEMO3/bulkprocess.aspx", form, ch)
}
func MakeJsonRequest(url string, form url.Values, ch chan string) {
	resp, _ := client.PostForm(url, form)

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	log.Println(resp.StatusCode)
	log.Println(url)
	log.Println(form)

	resp.Request.ParseMultipartForm(1024)
	f := resp.Request.Form

	log.Println(f)
	//result := new(Result)
	//json.NewDecoder(resp.Body).Decode(result)
	//secs := time.Since(start).Seconds()
	//
	//if len(result.MERCHREF) > 0 {
	//	log.Println(secs, " elapsed with redirect request ", " ", url)
	//}
}

var (
	frequency int
	debug     bool
)

func init() {
	RootCmd.AddCommand(crCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	crCmd.Flags().IntVarP(&frequency, "frequency", "f", 1, "how many times per second")
	crCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode")
}
