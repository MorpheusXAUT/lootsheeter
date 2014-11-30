// flags
package main

import (
	"flag"
)

var (
	debugLevelFlag    = flag.Int("debug", 3, "Sets the debug level (0-9), lower number displays more messages")
	httpPortFlag      = flag.Int("port", 3000, "Port for the webserver to bind to")
	httpHostFlag      = flag.String("host", "0.0.0.0", "Hostname for the webserver to bind to")
	mysqlUserFlag     = flag.String("mysqluser", "", "Username for authenticating to the MySQL server")
	mysqlPasswordFlag = flag.String("mysqlpassword", "", "Password for authenticating to the MySQL server")
	mysqlDatabaseFlag = flag.String("mysqldatabase", "", "Database to use with the MySQL server")
	mysqlHostFlag     = flag.String("mysqlhost", "localhost", "Hostname of the MySQL server")
	mysqlPortFlag     = flag.Int("mysqlport", 3306, "Port of the MySQL server")
)
