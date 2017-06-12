package cmd

import "log"

func fatal(msg interface{}) {
	log.Fatalln("Error: ", msg)
}
