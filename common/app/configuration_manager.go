package app

import "onevideogo/common/postgresql"

type ConfigurationManager struct {
	PostgreSqlConfig postgresql.Config
}

func NewConfigurationManager() *ConfigurationManager {
	postgreSqlConfig := getPostgreSQLConfig()
	return &ConfigurationManager{
		PostgreSqlConfig: postgreSqlConfig,
	}
}

func getPostgreSQLConfig() postgresql.Config {
	return postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		UserName:              "postgres",
		Password:              "postgres",
		DbName:                "productapp",
		MaxConnections:        "10",
		MaxConnectionIdleTime: "30s",
	}
}
