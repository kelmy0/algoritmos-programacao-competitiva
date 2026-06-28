package models

type Role struct {
	Id         string `db:"id"`
	Name       string `db:"name"`
	IsEmployee bool   `db:"is_employee"`
}
