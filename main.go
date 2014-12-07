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

	flag.Parse()

	if strings.EqualFold(*mysqlUserFlag, "") ||
		strings.EqualFold(*mysqlPasswordFlag, "") ||
		strings.EqualFold(*mysqlDatabaseFlag, "") ||
		strings.EqualFold(*ssoClientId, "") ||
		strings.EqualFold(*ssoClientSecret, "") {
		flag.Usage()
	}

	SetupLogger()

	InitialiseDatabase(*mysqlHostFlag, *mysqlPortFlag, *mysqlUserFlag, *mysqlPasswordFlag, *mysqlDatabaseFlag)

	InitialiseSessions()

	SetupRouter(true)

	HandleRequests(*httpHostFlag, *httpPortFlag)
}
