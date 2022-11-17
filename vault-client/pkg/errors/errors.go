package errors

import (
	"log"
	"fmt"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
}
