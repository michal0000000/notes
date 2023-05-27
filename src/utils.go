package utils

import "log"

func CheckErr(msg string, err error) {
	if err != nil {
		log.Printf(msg+": %s\n", err)
	}
}
