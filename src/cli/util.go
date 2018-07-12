package cli

import "log"

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
