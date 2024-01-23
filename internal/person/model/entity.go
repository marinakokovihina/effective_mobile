package model

import "time"

type PersonGender string

const (
	MaleGender   PersonGender = "male"
	FemaleGender PersonGender = "female"
	OtherGender  PersonGender = "other"
)

type Person struct {
	Id          int64         `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Surname     string        `json:"surname" db:"surname"`
	Patronymic  *string       `json:"patronymic" db:"patronymic"`
	Age         *int          `json:"age" db:"age"`
	Gender      *PersonGender `json:"gender" db:"gender"`
	Nationality *string       `json:"nationality" db:"nationality"`
	Created     time.Time     `json:"created" db:"created"`
	Updated     time.Time     `json:"updated" db:"updated"`
	Removed     bool          `json:"-" db:"removed"`
}
