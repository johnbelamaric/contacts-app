package main

import (
	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	FirstName    string
	MiddleName   string
	LastName     string
	EmailAddress string
}
