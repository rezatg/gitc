package main

import (
	"log"
	"os"
)

func main() {
	if err := Commands.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
