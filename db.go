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
		name = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True", c.User, c.Password, c.Host, c.Port, c.Name)
	default:
		return nil, fmt.Errorf(DBTypeNotSupportedMsg, c.Type)
	}
	return sql.Open(c.Type, name)
}

type DBStmt struct {
	// Student
	StFindStudentById           *sql.Stmt
	StFindStudentPasswordById   *sql.Stmt
	StUpdateStudentPasswordById *sql.Stmt
	// Session
	StCreateSessionByCookie             *sql.Stmt
	StCreateSessionByCookieWithLoggedIn *sql.Stmt
	StDeleteSessionById                 *sql.Stmt
	StDeleteSessionExpired              *sql.Stmt
	StVerifySessionNotExpired           *sql.Stmt
	StVerifySessionLoggedIn             *sql.Stmt
}

func (c *DBConf) Prepare(db *sql.DB) (*DBStmt, error) {
	var stmt DBStmt
	var err error
	switch c.Type {
	case "mysql":
		// Student
		stmt.StFindStudentById, err = db.Prepare(
			`SELECT s_id,s_name,s_mail,s_pass,s_c_id FROM student WHERE s_id=?`)
		if err != nil {
			return nil, err
		}
		stmt.StFindStudentPasswordById, err = db.Prepare(
			`SELECT s_pass FROM student WHERE s_id=?`)
		if err != nil {
			return nil, err
		}
		stmt.StUpdateStudentPasswordById, err = db.Prepare(
			`UPDATE student SET s_pass=? WHERE s_id=?`)
		if err != nil {
			return nil, err
		}
		// Session
		stmt.StCreateSessionByCookie, err = db.Prepare(
			`INSERT INTO session(ss_id,ss_expire) VALUES(?,?)`)
		if err != nil {
			return nil, err
		}
		stmt.StCreateSessionByCookieWithLoggedIn, err = db.Prepare(
			`INSERT INTO session(ss_id,ss_s_id,ss_expire) VALUES(?,?,?)`)
		if err != nil {
			return nil, err
		}
		stmt.StDeleteSessionById, err = db.Prepare(
			`DELETE FROM session WHERE ss_id=?`)
		if err != nil {
			return nil, err
		}
		stmt.StDeleteSessionExpired, err = db.Prepare(
			`DELETE FROM session WHERE ss_expire<NOW()`)
		if err != nil {
			return nil, err
		}
		stmt.StVerifySessionNotExpired, err = db.Prepare(
			`SELECT ss_expire>NOW() FROM session WHERE ss_id=?`)
		if err != nil {
			return nil, err
		}
		stmt.StVerifySessionLoggedIn, err = db.Prepare(
			`SELECT NOT ISNULL(ss_s_id) FROM session WHERE ss_id=?`)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(DBTypeNotSupportedMsg, c.Type)
	}
	return &stmt, nil
}
