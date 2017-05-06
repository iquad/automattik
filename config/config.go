package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// Config is a datatype to keep Configurations.
type Config struct {
	Host     string `json: "host"`
	Port     int    `json: "port"`
	DBType   string `json: "db_type"`
	DBParams string `json: "db_params"`
}

// Defaults defining default values for Config struct.
var Defaults = Config{"localhost", 9999, "sqlite3", "./automattik.db?cache_size=50"}

var allowedDatabaseTypes = map[string]bool{
	"sqlite3":  true,
	"postgres": true,
	"mysql":    true,
	"mssql":    true,
}

// NewConfig returns a new configuration file.
func NewConfig() *Config {
	var config Config
	config.Host = Defaults.Host
	config.Port = Defaults.Port
	config.DBType = Defaults.DBType
	config.DBParams = Defaults.DBParams
	return &config
}

// ProcessFlags is a defined type for flags and errors.
type ProcessFlags func() error

// BindFlags is a method for Config struct to binding flags for configuration.
func (config *Config) BindFlags() ProcessFlags {
	confFile := flag.String("conf", "", "path to the configuration file")
	dbType := flag.String("dbtype", Defaults.DBType, "database software")
	host := flag.String("host", Defaults.Host, "binding address of the server.")
	port := flag.Int("port", Defaults.Port, "port number of the database server.")
	dbParams := flag.String("dbparams", Defaults.DBParams, "parameters to open the database")

	return func() error {
		config.Host = *host
		config.Port = *port
		config.DBParams = *dbParams
		err := config.SetDBType(*dbType)
		if err != nil {
			return err
		}
		err = config.HandleConfFileFlag(*confFile)
		return err
	}
}

// HandleConfFileFlag handles configuration files
func (config *Config) HandleConfFileFlag(path string) error {
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return errors.New(fmt.Sprintf("Can't read file: %s", path))
		}

		err = config.Read(bufio.NewReader(file))
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to parse file: '%s' (%s).", path, err))
		}
	}
	return nil
}

// SetDBType sets database type.
func (config *Config) SetDBType(db_type string) error {
	if !allowedDatabaseTypes[db_type] {
		return errors.New(fmt.Sprintf("Unknown database type: '%s.", db_type))
	}

	config.DBType = db_type
	return nil
}

func (config *Config) Read(input io.Reader) error {
	return json.NewDecoder(input).Decode(config)
}

func (config *Config) Write(output io.Writer) error {
	return json.NewEncoder(output).Encode(config)
}

// Pretty pretties json
func (config *Config) Pretty(output io.Writer) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	data = append(data, []byte("\n")...)
	_, err = output.Write(data)
	return err
}
