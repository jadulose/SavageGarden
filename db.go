package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Type     string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Path     string
}

func (c *DBConfig) Open() (*sql.DB, error) {
	var name string
	switch c.Type {
	case "mysql":
		name = fmt.Sprintf("%s:%s@(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Name)
	default:
		return nil, fmt.Errorf("不支持的数据库类型：%s（支持mysql）", c.Type)
	}
	return sql.Open(c.Type, name)
}
