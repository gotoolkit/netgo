// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/spf13/cobra"
)

// freesslCmd represents the freessl command
var freesslCmd = &cobra.Command{
	Use:   "freessl",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var httpsSrv *http.Server

		// when testing locally it doesn't make sense to start
		// HTTPS server, so only do it in production.
		// In real code, I control this with -production cmd-line flag
		if inProduction {
			// Note: use a sensible value for data directory
			// this is where cached certificates are stored
			dataDir := "."
			hostPolicy := func(ctx context.Context, host string) error {
				// Note: change to your real domain
				allowedHost := "www.mydomain.com"
				if host == allowedHost {
					return nil
				}
				return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
			}

			httpsSrv = makeHTTPServer()
			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: hostPolicy,
				Cache:      autocert.DirCache(dataDir),
			}
			httpsSrv.Addr = ":443"
			httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

			go func() {
				err := httpsSrv.ListenAndServeTLS("", "")
				if err != nil {
					log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
				}
			}()
		}

		httpSrv := makeHTTPServer()
		httpSrv.Addr = ":80"
		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	},
}

const (
	htmlIndex    = `<html><body>Welcome!</body></html>`
	inProduction = true
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlIndex)
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleIndex)

	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func init() {
	RootCmd.AddCommand(freesslCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// freesslCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// freesslCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
