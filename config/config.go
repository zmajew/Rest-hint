package config

import (
	"flag"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Config struct {
	Enabled      bool
	DatabasePath string
	Port         string
	UserName     string
	Password     string
}

func NewConfig() *Config {
	databasePath := flag.String("bolt", "bolt://localhost:7687", "bolt adress of database (default: bolt://localhost:7687)")
	port := flag.Int("port", 9090, "port number (default 9090)")
	userName := flag.String("user", "neo4j", "")
	password := flag.String("password", "", "us")
	flag.Parse()
	return &Config{
		Enabled:      true,
		DatabasePath: *databasePath,
		Port:         strconv.Itoa(*port),
		UserName:     *userName,
		Password:     *password,
	}
}

func ConnectDatabase(config *Config) (neo4j.Driver, error) {
	return neo4j.NewDriver(config.DatabasePath, neo4j.BasicAuth(string(config.UserName), string(config.Password), ""))
}
