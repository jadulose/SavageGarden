package main

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
