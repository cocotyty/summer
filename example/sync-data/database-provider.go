package main

import (
	"database/sql"
	"github.com/cocotyty/summer"
)

func init() {
	summer.Put(&DatabaseProvider{})
}

type DatabaseProvider struct {
	DB *sql.DB
}

func (provider *DatabaseProvider) Init() {
	conn, err := sql.Open("mysql", Conf.MysqlDSN)
	if err != nil {
		panic(err)
	}
	provider.DB = conn
}

func (provider *DatabaseProvider) Provide() (db *sql.DB) {
	return provider.DB
}
