package utils

import (
	"fmt"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
