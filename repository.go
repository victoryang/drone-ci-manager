package main

import (
	"os"
	"path"

	"github.com/jinzhu/gorm"
)

func init() {
	ORM.AutoMigrate(&Repository{})
}

type Repository struct {
	gorm.Model      `json:"-"`

	Name 			string 		`gorm:"column:name;size:64"`
	GitUrl			string 		`gorm:"column:git_url;size:64"`
	RollingProject 		[]string 	`gorm:"column:rolling_project;size:64"`
	IsActived			bool 		`gorm:"column:is_actived;size:64"`
}

func NewRepository(name string, gitSshUrl string) error {
	proj := &Repository{
		Name: name,
		GitUrl: gitSshUrl,
		IsActived: true,
	}

	if err := ORM.Create(proj).Error; err != nil {
		return err
	}

	return nil
}