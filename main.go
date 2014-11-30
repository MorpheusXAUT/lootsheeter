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

	InitialiseDatabase(*mysqlHostFlag, *mysqlPortFlag, *mysqlUserFlag, *mysqlPasswordFlag, *mysqlDatabaseFlag)

	SetupRouter(true)

	HandleRequests(*httpHostFlag, *httpPortFlag)
}
