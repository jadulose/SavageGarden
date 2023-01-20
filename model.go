package main

import (
	"net/http"
	"time"
)

type Student struct {
	Id       string
	Name     string
	Mail     string
	Password []byte
	Cid      int
}

func (db *DBStmt) FindStudentById(id string) (*Student, error) {
	var stu Student
	err := db.StFindStudentById.QueryRow(id).
		Scan(&stu.Id, &stu.Name, &stu.Mail, &stu.Password, &stu.Cid)
	return &stu, err
}

func (db *DBStmt) FindStudentPasswordById(id string) ([]byte, error) {
	var pass []byte
	err := db.StFindStudentPasswordById.QueryRow(id).Scan(&pass)
	return pass, err
}

func (db *DBStmt) UpdateStudentPasswordById(id string, password []byte) error {
	_, err := db.StUpdateStudentPasswordById.Exec(password, id)
	return err
}

type Session struct {
	Id     string
	Sid    string
	Expire time.Time
}

func (db *DBStmt) CreateSessionByCookie(cookie *http.Cookie) error {
	_, err := db.StCreateSessionByCookie.Exec(cookie.Value, cookie.Expires)
	return err
}

func (db *DBStmt) CreateSessionByCookieWithLoggedIn(cookie *http.Cookie, sid string) error {
	_, err := db.StCreateSessionByCookieWithLoggedIn.Exec(cookie.Value, sid, cookie.Expires)
	return err
}

func (db *DBStmt) DeleteSessionById(id string) error {
	_, err := db.StDeleteSessionById.Exec(id)
	return err
}

func (db *DBStmt) DeleteSessionExpired() (int64, error) {
	rs, err := db.StDeleteSessionExpired.Exec()
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

func (db *DBStmt) VerifySessionNotExpired(id string) (bool, error) {
	var rs bool
	err := db.StVerifySessionNotExpired.QueryRow(id).Scan(&rs)
	return rs, err
}

func (db *DBStmt) VerifySessionLoggedIn(id string) bool {
	var rs bool
	err := db.StVerifySessionLoggedIn.QueryRow(id).Scan(&rs)
	return err == nil && rs
}
