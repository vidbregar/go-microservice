package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	d, err := strconv.Atoi(os.Args[1])
	if err != nil {
		os.Exit(1)
	}

	time.Sleep(time.Duration(d) * time.Second)
}
