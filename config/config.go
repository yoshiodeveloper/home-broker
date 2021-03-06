// Package config holds application configurations.
package config

import (
	"github.com/spf13/viper"
)

// PostgreSQLConfig holds the PostgreSQL configurations.
type PostgreSQLConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

// NewPostgreSQLConfigFromViper creates a new PostgreSQLConfig from viper.
func NewPostgreSQLConfigFromViper(v *viper.Viper) PostgreSQLConfig {
	c := PostgreSQLConfig{
		Host:     viper.GetString("DBHOST"),
		Port:     viper.GetInt("DBPORT"),
		Name:     viper.GetString("DBNAME"),
		User:     viper.GetString("DBUSER"),
		Password: viper.GetString("DBPASSWORD"),
	}
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.Name == "" {
		c.Name = "homebroker"
	}
	if c.User == "" {
		c.User = "homebroker"
	}
	if c.Password == "" {
		c.Password = "123456"
	}
	return c
}

// GinConfig holds the Gin server configurations.
type GinConfig struct {
	Port int
	Mode string // debug / release
}

// NewGinConfigFromViper creates a new DBGinConfig from viper.
func NewGinConfigFromViper(v *viper.Viper) GinConfig {
	c := GinConfig{
		Port: viper.GetInt("GINPORT"),
		Mode: viper.GetString("GINMODE"),
	}
	if c.Port == 0 {
		c.Port = viper.GetInt("PORT") // original gin var
		if c.Port == 0 {
			c.Port = 8080
		}
	}
	if c.Mode != "release" {
		c.Mode = viper.GetString("GIN_MODE") // original gin var
		if c.Mode != "release" {
			c.Mode = "debug"
		}
	}
	return c
}
