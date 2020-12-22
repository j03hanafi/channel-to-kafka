package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]

	msg, err := doProducer(broker, topic)

	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("\n\n\n %v", msg)

}
