package model

import "time"

type Person struct {
	ID        int
	Name      string
	BirthDate time.Time
}

func (p *Person) TableName() string {
	return "people"
}
