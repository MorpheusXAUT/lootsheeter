// flags
package main

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	DebugLevel      int
	HTTPPort        int
	HTTPHost        string
	MySqlUser       string
	MySqlPassword   string
	MySqlDatabase   string
	MySqlHost       string
	MySqlPort       int
	SSOClientID     string
	SSOClientSecret string
	SSOCallbackURL  string
}

var (
	config *Config
)

func ParseConfigFlags() (*Config, error) {
	debugLevelFlag := flag.Int("debug", 3, "Sets the debug level (0-9), lower number displays more messages")
	httpPortFlag := flag.Int("port", 3000, "Port for the webserver to bind to")
	httpHostFlag := flag.String("host", "0.0.0.0", "Hostname for the webserver to bind to")
	mysqlUserFlag := flag.String("mysqluser", "", "Username for authenticating to the MySQL server")
	mysqlPasswordFlag := flag.String("mysqlpassword", "", "Password for authenticating to the MySQL server")
	mysqlDatabaseFlag := flag.String("mysqldatabase", "", "Database to use with the MySQL server")
	mysqlHostFlag := flag.String("mysqlhost", "localhost", "Hostname of the MySQL server")
	mysqlPortFlag := flag.Int("mysqlport", 3306, "Port of the MySQL server")
	ssoClientIDFlag := flag.String("ssoid", "", "EVE Online Application Client ID")
	ssoClientSecretFlag := flag.String("ssosecret", "", "EVE Online Application Client Secret")
	ssoCallbackURLFlag := flag.String("ssocallback", "", "EVE Online Application Callback URL")
	configFileFlag := flag.String("config", "", "Config file to parse commandline parameters from")

	flag.Parse()

	var conf *Config

	if len(*configFileFlag) > 0 {
		configFile, err := os.Open(*configFileFlag)
		if err != nil {
			return &Config{}, err
		}

		decoder := json.NewDecoder(configFile)

		err = decoder.Decode(&conf)
		if err != nil {
			return &Config{}, err
		}
	} else {
		conf = &Config{
			DebugLevel:      *debugLevelFlag,
			HTTPPort:        *httpPortFlag,
			HTTPHost:        *httpHostFlag,
			MySqlUser:       *mysqlUserFlag,
			MySqlPassword:   *mysqlPasswordFlag,
			MySqlDatabase:   *mysqlDatabaseFlag,
			MySqlHost:       *mysqlHostFlag,
			MySqlPort:       *mysqlPortFlag,
			SSOClientID:     *ssoClientIDFlag,
			SSOClientSecret: *ssoClientSecretFlag,
			SSOCallbackURL:  *ssoCallbackURLFlag,
		}
	}

	return conf, nil
}
