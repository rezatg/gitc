package main

import (
	"log"
	"os"

	"github.com/rezatg/gitc/cmd"
)

func main() {
	if err := cmd.Commands.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
