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

	if strings.EqualFold(*mysqlUserFlag, "") || strings.EqualFold(*mysqlPasswordFlag, "") || strings.EqualFold(*mysqlDatabaseFlag, "") {
		flag.Usage()
	}

	SetupLogger()

	err := InitialiseDatabase(*mysqlHostFlag, *mysqlPortFlag, *mysqlUserFlag, *mysqlPasswordFlag, *mysqlDatabaseFlag)
	if err != nil {
		logger.Fatalf("Received error while initialising database: [%v]", err)
	}

	SetupRouter(true)

	HandleRequests(*httpHostFlag, *httpPortFlag)
}
