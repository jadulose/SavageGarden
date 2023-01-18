package main

import "database/sql"

type Student struct {
	Id       string
	Name     string
	Mail     string
	Password []byte
	Cid      int
}

func FindStudentById(db *sql.DB, id string) (*Student, error) {
	var stu Student
	err := db.QueryRow(`SELECT s_id,s_name,s_mail,s_pass,s_c_id FROM student WHERE s_id=?`, id).
		Scan(&stu.Id, &stu.Name, &stu.Mail, &stu.Password, &stu.Cid)
	return &stu, err
}

func UpdateStudentPasswordById(db *sql.DB, id string, password []byte) error {
	_, err := db.Exec(`UPDATE student SET s_pass=? WHERE s_id=?`, password, id)
	return err
}
