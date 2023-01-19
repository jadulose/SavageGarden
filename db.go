package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConf struct {
	Type     string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Path     string
}

const DBTypeNotSupportedMsg = "不支持的数据库类型：%s（目前支持mysql）"

func (c *DBConf) Open() (*sql.DB, error) {
	var name string
	switch c.Type {
	case "mysql":
		name = fmt.Sprintf("%s:%s@(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Name)
	default:
		return nil, fmt.Errorf(DBTypeNotSupportedMsg, c.Type)
	}
	return sql.Open(c.Type, name)
}

type DBStmt struct {
	// Student
	StmtFindStudentById           *sql.Stmt
	StmtUpdateStudentPasswordById *sql.Stmt
}

func (c *DBConf) Prepare(db *sql.DB) (*DBStmt, error) {
	var stmt DBStmt
	var err error
	switch c.Type {
	case "mysql":
		stmt.StmtFindStudentById, err = db.Prepare(
			`SELECT s_id,s_name,s_mail,s_pass,s_c_id FROM student WHERE s_id=?`)
		if err != nil {
			return nil, err
		}
		stmt.StmtUpdateStudentPasswordById, err = db.Prepare(
			`UPDATE student SET s_pass=? WHERE s_id=?`)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(DBTypeNotSupportedMsg, c.Type)
	}
	return &stmt, nil
}
