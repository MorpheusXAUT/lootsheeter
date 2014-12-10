// flags
package main

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	DebugLevel      int
	HttpPort        int
	HttpHost        string
	MySqlUser       string
	MySqlPassword   string
	MySqlDatabase   string
	MySqlHost       string
	MySqlPort       int
	SSOClientId     string
	SSOClientSecret string
	SSOCallbackUrl  string
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
	ssoClientIdFlag := flag.String("ssoid", "", "EVE Online Application Client ID")
	ssoClientSecretFlag := flag.String("ssosecret", "", "EVE Online Application Client Secret")
	ssoCallbackUrlFlag := flag.String("ssocallback", "", "EVE Online Application Callback URL")
	configFileFlag := flag.String("config", "", "Config file to parse commandline parameters from")

	flag.Parse()

	if len(*configFileFlag) > 0 {
		configFile, err := os.Open(*configFileFlag)
		if err != nil {
			return &Config{}, err
		}

		decoder := json.NewDecoder(configFile)

		var conf *Config

		err = decoder.Decode(&conf)
		if err != nil {
			return &Config{}, err
		}

		return conf, nil
	} else {
		conf := &Config{
			DebugLevel:      *debugLevelFlag,
			HttpPort:        *httpPortFlag,
			HttpHost:        *httpHostFlag,
			MySqlUser:       *mysqlUserFlag,
			MySqlPassword:   *mysqlPasswordFlag,
			MySqlDatabase:   *mysqlDatabaseFlag,
			MySqlHost:       *mysqlHostFlag,
			MySqlPort:       *mysqlPortFlag,
			SSOClientId:     *ssoClientIdFlag,
			SSOClientSecret: *ssoClientSecretFlag,
			SSOCallbackUrl:  *ssoCallbackUrlFlag,
		}

		return conf, nil
	}
}
