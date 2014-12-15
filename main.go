// lootsheeter project main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: lootsheeter [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	c, err := ParseConfigFlags()
	if err != nil {
		fmt.Printf("Received error while parsing config flags: [%v]\n", err)
		return
	}

	config = c

	if strings.EqualFold(config.MySqlUser, "") ||
		strings.EqualFold(config.MySqlPassword, "") ||
		strings.EqualFold(config.MySqlDatabase, "") ||
		strings.EqualFold(config.SSOClientID, "") ||
		strings.EqualFold(config.SSOClientSecret, "") ||
		strings.EqualFold(config.SSOCallbackURL, "") {
		flag.Usage()
	}

	SetupLogger()

	InitialiseDatabase()

	InitialiseSessions()

	SetupRouter()

	HandleRequests()
}
