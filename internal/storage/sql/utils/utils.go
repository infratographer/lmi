package utils

import (
	"database/sql"
	"fmt"

	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"
)

func dbConfigFromArgs(v *viper.Viper, dbName string) crdbx.Config {
	cfg := crdbx.Config{
		Name:     dbName,
		Host:     v.GetString("crdb.host"),
		User:     v.GetString("crdb.user"),
		Password: v.GetString("crdb.password"),
		Params:   v.GetString("crdb.params"),
		URI:      v.GetString("crdb.uri"),
	}

	cfg.Connections.MaxOpen = v.GetInt("crdb.connections.max_open")
	cfg.Connections.MaxIdle = v.GetInt("crdb.connections.max_idle")
	cfg.Connections.MaxLifetime = v.GetDuration("crdb.connections.max_lifetime")

	return cfg
}

func GetDBConnection(v *viper.Viper, dbName string, tracing bool) (*sql.DB, error) {
	cfg := dbConfigFromArgs(v, dbName)
	fmt.Printf("viper: %+v", v)
	fmt.Printf("cfg: %+v", cfg)
	return crdbx.NewDB(cfg, tracing)
}
