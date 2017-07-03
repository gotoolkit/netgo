package cmd

import "log"

func fatal(msg interface{}) {
	log.Fatalln("Error: ", msg)
}

func checkError(err error) {
	if err != nil {
		fatal(err.Error())
	}
}
